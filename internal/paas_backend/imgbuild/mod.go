package imgbuild

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/config"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/go-connections/nat"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type BuildError struct {
	Logs          string
	BuildErrorMsg string
}

func (e *BuildError) Error() string {
	return fmt.Sprintf("Build failed (error: %s)", e.BuildErrorMsg)
}

// Credits: https://github.com/go-git/go-git/issues/231#issuecomment-782835827
// Generates a tarball from the given branch of the repository
func writeTarFromCommit(repo *git.Repository, w io.Writer, hash *plumbing.Hash) error {
	// Get the corresponding commit hash.
	obj, err := repo.CommitObject(*hash)
	if err != nil {
		return err
	}

	// Let's have a look at the tree at that commit.
	tree, err := repo.TreeObject(obj.TreeHash)
	if err != nil {
		return err
	}

	type carrier struct {
		f *object.File
		r io.ReadCloser
	}

	files := make(chan carrier, 1000)
	g := &errgroup.Group{}

	g.Go(func() error {
		tarball := tar.NewWriter(w)

		for c := range files {
			err := tarball.WriteHeader(&tar.Header{
				Name: c.f.Name,
				Mode: 0600,
				Size: int64(c.f.Size),
			})
			if err != nil {
				return fmt.Errorf("failed to write header for file %s: %w", c.f.Name, err)
			}

			content, err := io.ReadAll(c.r)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", c.f.Name, err)
			}

			_, err = tarball.Write(content)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %w", c.f.Name, err)
			}

			err = c.r.Close()
			if err != nil {
				return fmt.Errorf("failed to close file %s: %w", c.f.Name, err)
			}
		}

		return tarball.Close()
	})

	addFile := func(f *object.File) error {
		fr, err := f.Reader()
		if err != nil {
			return err
		}

		files <- carrier{f, fr}

		return nil
	}

	err = tree.Files().ForEach(addFile)
	if err != nil {
		return err
	}

	close(files)

	return g.Wait()
}

// Builds an image from the last commit of a given branch, and assigns it the given tags
// On error, returns logs from the build process
// The branch worktree must contain a Dockerfile
func BuildGitBranch(state utils.State, app models.DBApplication, env models.DBEnvironment) (string, error) {
	repoPath := app.GetPath(state.Config)
	logrus.Debugf("Building image at %s", repoPath)

	buildCtx := bytes.Buffer{}
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open git repository: %w", err)
	}

	// Get commit hash of the branch
	rev := fmt.Sprintf("refs/remotes/origin/%s", env.Branch)
	hash, err := repo.ResolveRevision(plumbing.Revision(rev))
	if err != nil {
		return "", fmt.Errorf("failed to resolve revision %s: %w", rev, err)
	}

	// Generate image tag
	imageTag := fmt.Sprintf("%s/%s/%s:%s", config.Get().REGISTRY_REPO_URI, app.Name, env.Name, hash.String())

	// Generate tarball from the branch
	err = writeTarFromCommit(repo, &buildCtx, hash)
	if err != nil {
		return "", fmt.Errorf("failed to tar build context: %w", err)
	}

	// Build the image
	logrus.Debugf("Building image with tag %s", imageTag)
	buildOpts := types.ImageBuildOptions{
		Tags: []string{imageTag},
	}
	resp, err := state.DockerClient.ImageBuild(context.Background(), &buildCtx, buildOpts)
	if err != nil {
		return "", fmt.Errorf("failed to build image: %w", err)
	}
	defer resp.Body.Close()

	// Get messages from the build process
	buf := bytes.Buffer{}
	err = jsonmessage.DisplayJSONMessagesStream(resp.Body, &buf, 0, false, nil)
	if err != nil {
		return "", &BuildError{
			Logs:          buf.String(),
			BuildErrorMsg: err.Error(),
		}
	}

	logrus.Debugf("Built image at %s successfully", repoPath)
	return imageTag, nil
}

func GetExposedPort(dockerClient *client.Client, tag string) *int {
	imageInspect, _, err := dockerClient.ImageInspectWithRaw(context.Background(), tag)
	if err != nil {
		logrus.Errorf("Docker inspect error: %v", err)
		return nil
	}

	exposedPorts := imageInspect.Config.ExposedPorts
	if len(exposedPorts) == 1 {
		keys := make([]nat.Port, 0, len(exposedPorts))
		for p := range exposedPorts {
			keys = append(keys, p)
		}
		port := keys[0].Int()
		return &port
	} else if len(exposedPorts) > 1 {
		logrus.Warnf("Docker inspect warning - more than one exposed port in image")
	}
	return nil
}

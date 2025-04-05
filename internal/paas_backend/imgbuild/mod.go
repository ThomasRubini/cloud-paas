package imgbuild

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"

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
func writeTarFromBranch(repo *git.Repository, w io.Writer, revision string) error {
	hash, err := repo.ResolveRevision(plumbing.Revision(revision))
	if err != nil {
		return err
	}

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
func BuildGitBranch(repoPath string, branch string, tag string) error {
	logrus.Debugf("Building image at %s", repoPath)
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		logrus.Fatalf("cli error - %s", err)
	}

	buildOpts := types.ImageBuildOptions{
		Tags: []string{tag},
	}

	buildCtx := bytes.Buffer{}
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open git repository - %s", err)
	}

	err = writeTarFromBranch(repo, &buildCtx, fmt.Sprintf("refs/remotes/origin/%s", branch))
	if err != nil {
		return fmt.Errorf("failed to tar build context - %w", err)
	}

	resp, err := cli.ImageBuild(ctx, &buildCtx, buildOpts)
	if err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}
	defer resp.Body.Close()

	// Get messages from the build process
	buf := bytes.Buffer{}
	err = jsonmessage.DisplayJSONMessagesStream(resp.Body, &buf, 0, false, nil)
	if err != nil {
		return &BuildError{
			Logs:          buf.String(),
			BuildErrorMsg: err.Error(),
		}
	}

	logrus.Debugf("Built image at %s successfully", repoPath)
	return nil
}

func GetExposedPort(tag string) *int {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		logrus.Errorf("Docker client error - %s", err)
		return nil
	}
	imageInspect, _, err := cli.ImageInspectWithRaw(ctx, tag)
	if err != nil {
		logrus.Errorf("Docker inspect error - %s", err)
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

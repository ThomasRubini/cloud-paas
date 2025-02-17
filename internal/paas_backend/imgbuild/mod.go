package imgbuild

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
)

type BuildError struct {
	Logs          string
	BuildErrorMsg string
}

func (e *BuildError) Error() string {
	return fmt.Sprintf("Build failed (error: %s)", e.BuildErrorMsg)
}

// Builds an image from a directory containing a Dockerfile, and assigns it the given tags
// On error, returns logs from the build process
func Build(buildContextPath string, tags []string) error {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		log.Fatalf("cli error - %s", err)
	}

	buildOpts := types.ImageBuildOptions{
		Tags: tags,
	}

	//Archivage par la CLI ou le serveur ???
	buildCtx, err := archive.TarWithOptions(buildContextPath, &archive.TarOptions{})
	if err != nil {
		return fmt.Errorf("failed to tar build context - %s", err)
	}

	resp, err := cli.ImageBuild(ctx, buildCtx, buildOpts)
	if err != nil {
		return fmt.Errorf("failed to build image - %s", err)
	}
	defer resp.Body.Close()

	// récupère le flux du build docker (le print dans le terminal)
	buf := bytes.Buffer{}
	//c'est lui qui récupères jcrois et il écrit dans buf
	err = jsonmessage.DisplayJSONMessagesStream(resp.Body, &buf, 0, false, nil)
	if err != nil {
		return &BuildError{
			Logs:          buf.String(),
			BuildErrorMsg: err.Error(),
		}
	}

	return nil
}

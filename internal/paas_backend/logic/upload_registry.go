package logic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/config"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
	"github.com/docker/cli/cli/config/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/sirupsen/logrus"
)

func getAuth(conf *config.Config) string {
	auth := types.AuthConfig{
		Username: conf.REGISTRY_USER,
		Password: conf.REGISTRY_PASSWORD,
	}
	authStr, err := json.Marshal(auth)
	if err != nil {
		logrus.Errorf("Error marshalling auth: %v", err)
		return ""
	}

	return base64.StdEncoding.EncodeToString(authStr)
}

func UploadToRegistry(state utils.State, imageTag string) error {
	logrus.Debugf("Uploading image %v to registry..", imageTag)

	// Push the image to the registry
	resp, err := state.DockerClient.ImagePush(context.Background(), imageTag, image.PushOptions{
		RegistryAuth: getAuth(state.Config),
	})
	if err != nil {
		return fmt.Errorf("failed to push image - %w", err)
	}
	defer resp.Close()

	// Read output to verify there are no errors
	dec := json.NewDecoder(resp)
	for {
		var msg jsonmessage.JSONMessage
		err := dec.Decode(&msg)
		if err != nil {
			break
		}
		if msg.Error != nil {
			return fmt.Errorf("push error - %s", msg.Error.Message)
		}
	}

	logrus.Debugf("Push response is successful")
	return nil
}

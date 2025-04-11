package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/ThomasRubini/cloud-paas/internal/paas_cli/config"
	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v2"
)

func GetAPIClient() *resty.Client {
	r := resty.New()
	r.SetBaseURL(config.Get().BACKEND_URL)
	return r
}

func OpenInEditor(filePath string) ([]byte, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		fallbackEditors := []string{"micro", "nano", "vim", "vi"}
		for _, fallbackEditor := range fallbackEditors {
			if _, err := exec.LookPath(fallbackEditor); err == nil {
				editor = fallbackEditor
				break
			}
		}
	}
	if editor == "" {
		return nil, fmt.Errorf("no editor found, please install one")
	}

	openEditorCmd := exec.Command(editor, filePath)
	openEditorCmd.Stdout = os.Stdout
	openEditorCmd.Stderr = os.Stderr
	openEditorCmd.Stdin = os.Stdin

	if err := openEditorCmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to open editor: %w", err)
	}

	// Read the updated environment variables from the temp file
	updatedFile, err := os.ReadFile(filePath) // Use filePath instead of tempFile.Name()
	if err != nil {
		return nil, fmt.Errorf("failed to read temp file: %w", err)
	}

	return updatedFile, nil
}

func YAMLtoJSON(yamlContent []byte) ([]byte, error) {
	parsedContent := make(map[string]interface{})
	err := yaml.Unmarshal(yamlContent, &parsedContent)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling content: %w", err)
	}

	if len(parsedContent) == 0 {
		return []byte{}, nil
	}

	for _, v := range parsedContent {
		if _, ok := v.(map[interface{}]interface{}); ok {
			return nil, fmt.Errorf("content contains nested maps, which are not supported")
		}
	}

	// Marshal the struct into JSON
	jsonData, err := json.Marshal(parsedContent)
	if err != nil {
		return nil, fmt.Errorf("error marshalling to JSON: %w", err)
	}
	return jsonData, nil
}

// JSONtoYAML converts JSON data to YAML format
func JSONtoYAML(jsonData []byte) ([]byte, error) {
	if len(jsonData) == 0 {
		return []byte{}, nil
	}

	var data interface{}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Marshal the interface{} into YAML
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to YAML: %w", err)
	}

	return yamlData, nil
}

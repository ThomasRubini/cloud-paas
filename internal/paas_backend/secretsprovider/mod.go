package secretsprovider

import (
	"encoding/json"
	"fmt"
	"os"
)

type SecretsProvider interface {
	// returns an empty string if the key doesnt have a secret associated
	GetSecret(key string) (string, error)
	SetSecret(key, value string) error
}

type FileSecretsProvider struct {
	filePath string
}

func FromFile(filePath string) *FileSecretsProvider {
	return &FileSecretsProvider{filePath: filePath}
}

func (fsp *FileSecretsProvider) readData() (map[string]string, error) {
	data, err := os.ReadFile(fsp.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}

	secrets := make(map[string]string)
	err = json.Unmarshal(data, &secrets)
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

func (fsp *FileSecretsProvider) SetSecret(key, value string) error {
	secrets, err := fsp.readData()
	if err != nil {
		return fmt.Errorf("could not read secrets file: %v", err)
	}

	secrets[key] = value

	data, err := json.Marshal(secrets)
	if err != nil {
		return err
	}

	err = os.WriteFile(fsp.filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (fsp *FileSecretsProvider) GetSecret(key string) (string, error) {
	secrets, err := fsp.readData()
	if err != nil {
		return "", fmt.Errorf("could not read secrets file: %v", err)
	}

	secret, ok := secrets[key]
	if !ok {
		return "", nil
	}

	return secret, nil
}

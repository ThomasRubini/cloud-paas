package secretsprovider

import (
	"encoding/json"
	"fmt"
	"os"
)

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

func (fsp *FileSecretsProvider) writeData(secrets map[string]string) error {
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
func (fsp *FileSecretsProvider) SetSecret(key, value string) error {
	secrets, err := fsp.readData()
	if err != nil {
		return fmt.Errorf("could not read secrets file: %w", err)
	}

	secrets[key] = value

	if fsp.writeData(secrets) != nil {
		return fmt.Errorf("could not write secrets file: %w", err)
	}

	return nil
}

func (fsp *FileSecretsProvider) GetSecret(key string) (string, error) {
	secrets, err := fsp.readData()
	if err != nil {
		return "", fmt.Errorf("could not read secrets file: %w", err)
	}

	secret, ok := secrets[key]
	if !ok {
		return "", nil
	}

	return secret, nil
}

func (fsp *FileSecretsProvider) DeleteSecret(key string) error {
	secrets, err := fsp.readData()
	if err != nil {
		return fmt.Errorf("could not read secrets file: %w", err)
	}

	delete(secrets, key)

	if fsp.writeData(secrets) != nil {
		return fmt.Errorf("could not write secrets file: %w", err)
	}

	return nil
}

package secretsprovider

type Core interface {
	// returns an empty string if the key doesnt have a secret associated
	GetSecret(key string) (string, error)
	SetSecret(key, value string) error
	DeleteSecret(key string) error
}

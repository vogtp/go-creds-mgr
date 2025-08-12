package filestorage

type Option func(*fileStorage)

func SecretsPassword(password string) Option {
	return func(ts *fileStorage) {
		ts.secretsPassword = []byte(password)
	}
}

func StorePath(path string) Option {
	return func(ts *fileStorage) {
		ts.storagePath = path
	}
}

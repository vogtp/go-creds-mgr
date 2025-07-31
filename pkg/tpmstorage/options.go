package tpmstorage

type Option func(*tpmStorage)

func SecretsPassword(password string) Option {
	return func(ts *tpmStorage) {
		ts.secretsPassword = []byte(password)
	}
}

func StorePath(path string) Option {
	return func(ts *tpmStorage) {
		ts.storagePath = path
	}
}

func TPMDevice(path ...string) Option {
	return func(ts *tpmStorage) {
		ts.tpmDevPaths = path
	}
}

package config

type AppConfig struct {
	HTTP     HTTPConfig
	Database DatabaseConfig
}

func Load() (AppConfig, error) {
	if err := LoadDotEnv(".env", ".env.local"); err != nil {
		return AppConfig{}, err
	}

	return AppConfig{
		HTTP:     loadHTTPConfigFromEnv(),
		Database: loadDatabaseConfigFromEnv(),
	}, nil
}

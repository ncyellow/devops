package config

// GeneralConfig содержит общие параметры по настройке агента и сервера
type GeneralConfig struct {
	Address   string `env:"ADDRESS"`
	SecretKey string `env:"KEY"`
}

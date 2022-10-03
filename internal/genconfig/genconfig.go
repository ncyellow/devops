// Package genconfig содержит общие параметры конфигурации, которые есть и на сервере и на агенте
package genconfig

// GeneralConfig содержит общие параметры по настройке агента и сервера
type GeneralConfig struct {
	Address   string `env:"ADDRESS"`
	SecretKey string `env:"KEY"`
	CryptoKey string `env:"CRYPTO_KEY"`
}

func (g *GeneralConfig) GeneralCfg() *GeneralConfig {
	return g
}

// Package genconfig содержит общие параметры конфигурации, которые есть и на сервере и на агенте
package genconfig

// GeneralConfig содержит общие параметры по настройке агента и сервера
type GeneralConfig struct {
	Address   string `env:"ADDRESS" json:"address"`
	SecretKey string `env:"KEY" json:"secret_key"`
	CryptoKey string `env:"CRYPTO_KEY" json:"crypto_key"`
}

func (g *GeneralConfig) GeneralCfg() *GeneralConfig {
	return g
}

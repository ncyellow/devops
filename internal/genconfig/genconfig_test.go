package genconfig

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneralConfig_GeneralCfg(t *testing.T) {
	conf := GeneralConfig{
		Address:   "address",
		SecretKey: "secret",
	}
	assert.Equal(t, conf.Address, "address")
	assert.Equal(t, conf.SecretKey, "secret")
	assert.Equal(t, conf.GeneralCfg(), &conf)
}

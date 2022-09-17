package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestConfig базовая проверка что конфиг агента в порядке
func TestConfig(t *testing.T) {
	conf := Config{
		ReportInterval: time.Second * 10,
		PollInterval:   time.Second * 20,
	}
	assert.Equal(t, conf.ReportInterval, time.Second*10)
	assert.Equal(t, conf.PollInterval, time.Second*20)
}

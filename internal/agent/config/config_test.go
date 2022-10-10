package config

import (
	"testing"
	"time"

	"github.com/ncyellow/devops/internal/genconfig"
	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want Config
	}{
		{
			"file not exists",
			args{
				fileName: "file not exists",
			},
			Config{},
		},
		{
			"broken file",
			args{
				fileName: "test_data/broken_agent_config.json",
			},
			Config{},
		},
		{
			"correct json",
			args{
				fileName: "test_data/agent_config.json",
			},
			Config{
				GeneralConfig: genconfig.GeneralConfig{
					Address:   "localhost:8080",
					CryptoKey: "/path/to/key.pem",
				},
				ReportInterval: genconfig.Duration{Duration: time.Second * 12},
				PollInterval:   genconfig.Duration{Duration: time.Second * 32},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadConfig(tt.args.fileName); !assert.Equal(t, got, tt.want) {
				t.Errorf("ReadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

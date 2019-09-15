package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

var (
	global *Config
)

// LoadGlobalConfig
func LoadGlobalConfig(fpath string) error {
	c, err := ParseConfig(fpath)
	if err != nil {
		return err
	}
	global = c
	return nil
}

// GetGlobalConfig
func GetGlobalConfig() *Config {
	if global == nil {
		return &Config{}
	}
	return global
}

// ParseConfig 解析配置文件
func ParseConfig(fpath string) (*Config, error) {
	var c Config
	_, err := toml.DecodeFile(fpath, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Config
type Config struct {
	RunMode string `toml:"run_mode"`
}

type PodTemplate struct {
	BaseImage     string `toml:"base_image"`
	CPURequest    string `toml:"cpu_request"`
	MemoryRequest string `toml:"memory_request"`
	IsGPU         bool   `toml:"is_gpu"`
}

package config

import (
	"errors"

	"github.com/redhat-et/copilot-ops/pkg/ai"
	"github.com/redhat-et/copilot-ops/pkg/ai/bloom"
	"github.com/redhat-et/copilot-ops/pkg/ai/gpt3"
	"github.com/redhat-et/copilot-ops/pkg/ai/gptj"
	"github.com/spf13/viper"
)

const (
	ConfigName      = ".copilot-ops"
	ConfigFile      = ".copilot-ops.yaml"
	ConfigFileLocal = ".copilot-ops.local"
)

// Config Defines the struct into which the config-file will be parsed.
type Config struct {
	Filesets []Filesets `json:"filesets,omitempty" yaml:"filesets,omitempty"`
	// OpenAI Defines the settings necessary for the OpenAI GPT-3 backend.
	// FIXME: rename to GPT-3
	OpenAI *gpt3.Config `json:"openAI,omitempty" yaml:"openAI,omitempty"`
	// Backend Defines which AI backend should be used in order to generate completions.
	// Valid models include: gpt-3, gpt-j, opt, and bloom.
	Backend ai.Backend `json:"backend"`
	// GPTJ Defines the configuration options for using GPT-J.
	GPTJ *gptj.Config `json:"gptj,omitempty" yaml:"gptj,omitempty"`
	// BLOOM Defines the configuration for using BLOOM.
	BLOOM *bloom.Config `json:"bloom,omitempty" yaml:"bloom,omitempty"`
}

type Filesets struct {
	Name  string   `json:"name" yaml:"name"`
	Files []string `json:"files" yaml:"files"`
}

// OpenAI Defines the settings for accessing and using OpenAI's tooling.
// GPTJ Defines the structure required for configuring GPT-J.
type GPTJ struct {
	URL string `json:"url,omitempty" yaml:"url,omitempty"`
}

// Load the config from file if it exists, but if it doesn't exist
// we'll just use the defaults and continue without error.
// Errors here might return if the file exists but is invalid.
func (c *Config) Load() error {
	// bind to environment variables
	openAIEnvs := map[string]string{
		"openai.apikey": "OPENAI_API_KEY",
		"openai.orgid":  "OPENAI_ORG_ID",
		"openai.url":    "OPENAI_URL",
	}
	for k, v := range openAIEnvs {
		if err := viper.BindEnv(k, v); err != nil {
			return err
		}
	}
	viper.SetEnvPrefix("COPILOT_OPS")
	viper.AutomaticEnv()

	// paths to look for the config file in
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("$HOME")
	// viper.AddConfigPath("..") // parent? grandparent? grandgrandparent?
	viper.AddConfigPath(".")

	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	viper.SetConfigName(ConfigName) // name of config file (without extension)

	if err := viper.MergeInConfig(); err != nil {
		var configFileNotFound viper.ConfigFileNotFoundError
		if ok := errors.As(err, &configFileNotFound); !ok {
			return err // allow no config file
		}
	}

	// optionally look for local (gitignored) config file and merge it in
	viper.SetConfigName(ConfigFileLocal)

	if err := viper.MergeInConfig(); err != nil {
		var configFileNotFound viper.ConfigFileNotFoundError
		if ok := errors.As(err, &configFileNotFound); !ok {
			return err // allow no config file
		}
	}

	if err := viper.Unmarshal(c); err != nil {
		return err
	}

	return nil
}

// FindFileset Returns a fileset with the matching name,
// or nil if none exists.
func (c *Config) FindFileset(name string) *Filesets {
	for _, fileset := range c.Filesets {
		if fileset.Name == name {
			return &fileset
		}
	}
	return nil
}

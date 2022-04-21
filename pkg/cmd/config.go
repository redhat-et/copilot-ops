package cmd

import (
	"log"

	"github.com/spf13/viper"
)

const CONFIG_FILE = ".copilot-ops.yaml"

type Config struct {
	Filesets []ConfigFilesets
	OpenAI   ConfigOpenAI
}

type ConfigFilesets struct {
	Name  string
	Files []string
}

type ConfigOpenAI struct {
	ApiKey string
	OrgId  string
}

// Load the config from file if it exists, but if it doesn't exist
// we'll just use the defaults and continue without error.
// Errors here might return if the file exists but is invalid.
func (c *Config) Load() error {

	// bind to environment variables
	viper.BindEnv("openai.apikey", "OPENAI_API_KEY")
	viper.BindEnv("openai.orgid", "OPENAI_ORG_ID")
	viper.SetEnvPrefix("COPILOT_OPS")
	viper.AutomaticEnv()

	// paths to look for the config file in
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("$HOME")
	// viper.AddConfigPath("..") // parent? grandparent? grandgrandparent?
	viper.AddConfigPath(".")

	viper.SetConfigType("yaml")         // REQUIRED if the config file does not have the extension in the name
	viper.SetConfigName(".copilot-ops") // name of config file (without extension)

	if err := viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err // allow no config file
		}
	}

	log.Printf("viper: %+v\n", viper.ConfigFileUsed())

	// optionally look for local (gitignored) config file and merge it in
	viper.SetConfigName(".copilot-ops.local")

	if err := viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err // allow no config file
		}
	}

	log.Printf("viper: %+v\n", viper.ConfigFileUsed())

	if err := viper.Unmarshal(c); err != nil {
		return err
	}

	return nil
}

func (c *Config) FindFileset(name string) *ConfigFilesets {
	for _, fileset := range c.Filesets {
		if fileset.Name == name {
			return &fileset
		}
	}
	return nil
}

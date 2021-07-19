package runner

import (
	"math/rand"
	"os"
	"time"

	"gopkg.in/yaml.v3"
	"iplookup/subscraping"
)

// MultipleKeyPartsLength is the max length for multiple keys
const MultipleKeyPartsLength = 2

// YAMLIndentCharLength number of chars for identation on write YAML to file
const YAMLIndentCharLength = 4

// ConfigFile contains the fields stored in the configuration file
type ConfigFile struct {
	// Sources contains a list of sources to use for enumeration
	Sources []string `yaml:"sources,omitempty"`
	// AllSources contains the list of all sources for enumeration (slow)
	AllSources []string `yaml:"all-sources,omitempty"`
	// ExcludeSources contains the sources to not include in the enumeration process
	ExcludeSources []string `yaml:"exclude-sources,omitempty"`

	Proxy string `yaml:"proxy"`
	// API keys for different sources  # TODO 默认KEY列表
	Dnsgrep []string `yaml:"dnsgrep"`
	C99     []string `yaml:"c99"`
	// Version indicates the version of subfinder installed.
	Version string `yaml:"iplookup-version"`
}

// GetConfigDirectory gets the subfinder config directory for a user
func GetConfigDirectory() (string, error) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	var config string

	directory, err := os.UserHomeDir()
	if err != nil {
		return config, err
	}
	config = directory + "/.config/iplookup"

	// Create All directory for subfinder even if they exist
	err = os.MkdirAll(config, os.ModePerm)
	if err != nil {
		return config, err
	}

	return config, nil
}

// CheckConfigExists checks if the config file exists in the given path
func CheckConfigExists(configPath string) bool {
	if _, err := os.Stat(configPath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}

// MarshalWrite writes the marshaled yaml config to disk
func (c *ConfigFile) MarshalWrite(file string) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	// Indent the spaces too
	enc := yaml.NewEncoder(f)
	enc.SetIndent(YAMLIndentCharLength)
	err = enc.Encode(&c)
	f.Close()
	return err
}

// UnmarshalRead reads the unmarshalled config yaml file from disk
func UnmarshalRead(file string) (ConfigFile, error) {
	config := ConfigFile{}

	f, err := os.Open(file)
	if err != nil {
		return config, err
	}
	err = yaml.NewDecoder(f).Decode(&config)
	f.Close()
	return config, err
}

// GetKeys gets the API keys from config file and creates a Keys struct
// We use random selection of api keys from the list of keys supplied.
// Keys that require 2 options are separated by colon (:).
func (c *ConfigFile) GetKeys() subscraping.Keys {
	keys := subscraping.Keys{}

	if len(c.Dnsgrep) > 0 {
		keys.Dnsgrep = c.Dnsgrep[rand.Intn(len(c.Dnsgrep))]
	}
	if len(c.C99) > 0 {
		keys.C99 = c.C99[rand.Intn(len(c.C99))]
	}

	return keys
}

func (c *ConfigFile) GetProxy() subscraping.Proxy {
	var proxy subscraping.Proxy
	if len(c.Proxy) > 0 {
		proxy = subscraping.Proxy(c.Proxy)
	}

	return proxy
}

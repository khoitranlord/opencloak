package configs

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	custerror "github.com/khoitranlord/opencloak/src/internal/error"
	"gopkg.in/yaml.v3"
)

var once sync.Once

var globalConfigs *Configs

type Configs struct {
	Public   HttpConfigs     `json:"public,omitempty" yaml:"public,omitempty"`
	Private  HttpConfigs     `json:"private,omitempty" yaml:"private,omitempty"`
	Logger   LoggerConfigs   `json:"logger,omitempty" yaml:"logger,omitempty"`
	Database DatabaseConfigs `json:"database,omitempty" yaml:"sqlite,omitempty"`
}

func (c Configs) String() string {
	configBytes, _ := json.Marshal(c)
	return string(configBytes)
}

func Init(ctx context.Context) {
	once.Do(func() {
		configs, err := readConfig()
		if err != nil {
			log.Fatal(err)
			return
		}
		globalConfigs = configs
	})
}

func Get() *Configs {
	return globalConfigs
}

type HttpConfigs struct {
	Name string           `json:"name,omitempty" yaml:"name,omitempty"`
	Port int              `json:"port,omitempty" yaml:"port,omitempty"`
	Tls  TlsConfig        `json:"tls,omitempty" yaml:"tls,omitempty"`
	Auth BasicAuthConfigs `json:"auth,omitempty" yaml:"auth,omitempty"`
}

type TlsConfig struct {
	Cert      string `json:"cert,omitempty" yaml:"cert,omitempty"`
	Key       string `json:"key,omitempty" yaml:"key,omitempty"`
	Authority string `json:"authority,omitempty" yaml:"authority,omitempty"`
}

type DatabaseConfigs struct {
	Connection string `json:"connection" yaml:"connection,omitempty"`
}

func (c TlsConfig) Enabled() bool {

	if len(c.Cert) > 0 && len(c.Key) > 0 {
		return true
	}
	return false
}

type LoggerConfigs struct {
	Level    string `json:"level,omitempty" yaml:"level,omitempty"`
	Encoding string `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

type BasicAuthConfigs struct {
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	Token    string `json:"token,omitempty" yaml:"token,omitempty"`
}

func readConfig() (*Configs, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	configFile, err := readConfigFile(path)
	if err != nil {
		return nil, err
	}

	configs, err := parseConfig(configFile)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func getConfigFilePath() (string, error) {
	path := os.Getenv(ENV_CONFIG_FILE_PATH)
	if len(path) == 0 {
		return "", custerror.FormatNotFound("ENV_CONFIG_FILE_PATH not found, unable to read configurations")
	}
	return path, nil
}

func readConfigFile(path string) ([]byte, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, custerror.FormatInternalError("readConfigFile: err = %s", err)
	}

	return contents, nil
}

func parseConfig(contents []byte) (*Configs, error) {
	configs := &Configs{}
	if jsonErr := json.Unmarshal(contents, configs); jsonErr != nil {
		if yamlErr := yaml.Unmarshal(contents, configs); yamlErr != nil {
			return nil, custerror.FormatInvalidArgument("parseConfig: config parse JSON err = %s YAML err = %s", jsonErr, yamlErr)
		}
	}
	return configs, nil
}

package config

import (
	"easy-api-prom-alert-sms/logging"

	"flag"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type ConfigFlag struct {
	ConfigFile  string
	ListenPort  string
	EnableHttps bool
	CertFile    string
	KeyFile     string
}

func newConfigFlag(configFile string, listenPort string, enableHttps bool, certFile string, keyFile string) *ConfigFlag {
	return &ConfigFlag{configFile, listenPort, enableHttps, certFile, keyFile}
}

func LoadConfigFlag(validate *validator.Validate) (*ConfigFlag, error) {
	var (
		configFile  string
		listenPort  int
		enableHttps string
		certFile    string
		keyFile     string
	)

	flag.StringVar(&configFile, "config-file", "fixtures/config.default.yaml", "configuration file path")
	flag.StringVar(&certFile, "cert-file", "fixtures/tls/server.crt", "certificat file path")
	flag.StringVar(&keyFile, "key-file", "fixtures/tls/server.key", "private key file path")
	flag.StringVar(&enableHttps, "enable-https", "false", "configuration to enable https")
	flag.IntVar(&listenPort, "port", 5957, "listening port")
	flag.Parse()

	if err := validate.Var(listenPort, "required,min=1024,max=49151"); err != nil {
		logging.Log(logging.Error, "you should provide a port number between 1024 and 49151")
		return nil, err
	}

	boolEnableHttps, err := strconv.ParseBool(enableHttps)
	if err != nil {
		logging.Log(logging.Error, "unable to parse enable-https flag")
		return nil, err
	}

	return newConfigFlag(configFile, strconv.Itoa(listenPort), boolEnableHttps, certFile, keyFile), nil
}

package config

import (
	"flag"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type ConfigFlag struct {
	ConfigFile  string `validate:"required"`
	ListenPort  int    `validate:"required,gte=1024,lte=49151"`
	EnableHttps bool
	CertFile    string `validate:"required_if=EnableHttps true"`
	KeyFile     string `validate:"required_if=EnableHttps true"`
}

func newConfigFlag(configFile string, listenPort int, enableHttps bool, certFile string, keyFile string) *ConfigFlag {
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

	boolEnableHttps, err := strconv.ParseBool(enableHttps)
	if err != nil {
		return nil, err
	}

	cfgflag := newConfigFlag(configFile, listenPort, boolEnableHttps, certFile, keyFile)
	if err := validate.Struct(cfgflag); err != nil {
		return nil, err
	}

	return cfgflag, nil
}

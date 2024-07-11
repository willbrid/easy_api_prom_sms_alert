package file

import (
	"fmt"
	"os"
)

func CreateConfigFileForTesting(configContent string) (string, error) {
	configFile, err := os.CreateTemp("", "testconfig-*.yaml")
	if err != nil {
		return "", fmt.Errorf("unable to create temp file : %s", err.Error())
	}

	_, err = configFile.WriteString(configContent)
	if err != nil {
		return "", fmt.Errorf("unable to write to temp file : %s", err.Error())
	}
	configFile.Close()

	return configFile.Name(), nil
}

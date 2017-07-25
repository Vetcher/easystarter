package backend

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kpango/glg"
	"gopkg.in/ini.v1"
)

const ENV_SECTION = ""

var (
	ENV_FILE_NAME = *flag.String("config", "env.ini", "File with configuration parameters")

	Environment *ini.File
)

func CurrentEnvironmentString() string {
	formattedStr := []string{""}
	for _, key := range Environment.Section(ENV_SECTION).Keys() {
		formattedStr = append(formattedStr, fmt.Sprintf("%v=%v", key.Name(), key.Value()))
	}
	return strings.Join(formattedStr, "\n")
}

func AllEnvironmentString() string {
	cmd := exec.Command("printenv")
	s, err := cmd.Output()
	if err != nil {
		glg.Errorf("Can't print environment: %v", err)
		return ""
	}
	return string(s)
}

func SetupEnv() bool {
	var err error
	Environment, err = LoadEnv()
	if err != nil {
		glg.Warnf("Can't load environment: %v", err)
		return false
	}
	err = ExpandEnv(Environment)
	if err != nil {
		glg.Errorf("%v", err)
		return false
	}
	return true
}

func CreateEnvFile() {
	file, err := os.Create(ENV_FILE_NAME)
	if err != nil {
		glg.Warnf("Can't create file `%v`: %v.", ENV_FILE_NAME, err.Error())
	}
	glg.Infof("File `%v` was created.", ENV_FILE_NAME)
	glg.Info("You can add some environment variables.")
	file.Close()
}

func LoadEnv() (*ini.File, error) {
	file, err := os.Open(ENV_FILE_NAME)
	if err != nil {
		CreateEnvFile()
	} else {
		file.Close()
	}
	envConfig, err := ini.Load(ENV_FILE_NAME)
	if err != nil {
		return nil, err
	}
	return envConfig, nil
}

func ExpandEnv(cfg *ini.File) error {
	keys := cfg.Section("").Keys()
	var err error
	for _, key := range keys {
		env, ok := os.LookupEnv(key.Name())
		if ok {
			err = os.Setenv(key.Name(), fmt.Sprintf("%s:%s", env, os.ExpandEnv(key.Value())))
		} else {
			err = os.Setenv(key.Name(), os.ExpandEnv(key.Value()))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

package config

import "os"

import "gopkg.in/yaml.v3"

type letterboxd struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type emby struct {
	Username string `yaml:"username"`
}

type user struct {
	Letterboxd letterboxd `yaml:"letterboxd"`
	Emby       emby       `yaml:"emby"`
}

type Config struct {
	Users []user `yaml:"users"`
}

func Load(filename string) Config {
	var data, readErr = os.ReadFile(filename)
	if readErr != nil {
		panic(readErr)
	}

	var config Config
	if yamlErr := yaml.Unmarshal(data, &config); yamlErr != nil {
		panic(yamlErr)
	}
	return config
}

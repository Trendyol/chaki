package config

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func parseFile(path string) (string, string, string) {
	var (
		dir  = filepath.Dir(path)
		file = filepath.Base(path)
		ext  = filepath.Ext(file)
	)
	return dir, strings.TrimSuffix(file, ext), ext[1:]
}

func readConfig(path string) (*viper.Viper, error) {
	dir, file, ext := parseFile(path)
	v := viper.New()
	v.AddConfigPath(dir)
	v.SetConfigName(file)
	v.SetConfigType(ext)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v, nil
}

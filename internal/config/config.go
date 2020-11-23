package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	configExt string = ".yaml"
)

var (
	errConfigNotFound = errors.New("Config: config not found")
	otherDirs         = []string{"/etc"}
)

func defaultConfigFile() (string, error) {
	cf := filepath.Base(os.Args[0])
	fmt.Println(cf)
	bd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	n := strings.LastIndexByte(cf, '.')
	if n > 0 {
		cf = cf[:n]
	}
	cf += configExt
	if _, err := os.Stat(cf); !os.IsNotExist(err) {
		return cf, nil
	}
	bd = filepath.Join(bd, cf)
	if _, err := os.Stat(bd); !os.IsNotExist(err) {
		return bd, nil
	}
	for _, v := range otherDirs {
		bd = filepath.Join(v, cf)
		if _, err := os.Stat(bd); !os.IsNotExist(err) {
			return bd, nil
		}
	}
	return cf, errConfigNotFound
}

// GetConfig returns content of config file
func GetConfig() ([]byte, error) {
	cf, err := defaultConfigFile()
	if err != nil {
		fmt.Println("h", err)
		return nil, err
	}
	cont, err := ioutil.ReadFile(cf)
	if err != nil {
		fmt.Println("t", err)
		return nil, err
	}
	return cont, nil
}

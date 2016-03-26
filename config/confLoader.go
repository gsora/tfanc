package config

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/user"
	"sort"

	"github.com/gsora/tfanc/config"
)

// WARNING: never directly use this variable before calling getUserHome()
var configFilePath = "/.config/tfanc.json"

func getUserHome() error {
	cU, err := user.Current()
	if err != nil {
		return err
	}

	configFilePath = cU.HomeDir + configFilePath

	return nil
}

// LoadConfig loads a configuration file from the standard path, defined by "configFilePath"
func LoadConfig() (config.Configuration, error) {

	if err := getUserHome(); err != nil {
		return err
	}

	f, _ := os.Open(configFilePath)
	defer f.Close()

	fReader := bufio.NewReader(f)

	buf := new(bytes.Buffer)
	buf.ReadFrom(fReader)

	var m config.Configuration
	json.Unmarshal(buf.Bytes(), &m)

	//
	// Configuration file sanity checks
	//

	// The first condition is true only if "targets" in the json file is not even defined, so
	// we'll check for array length too.
	if m.Targets == nil || len(m.Targets) == 0 {
		return config.Configuration{}, errors.New("configuration error: no \"targets\" field defined, cannot continue without. Check your configuration file, it's malformed")
	}

	// Assuming from here that our configuration file contains at least one target, we need to correctly
	// parse what the user was thinking when writing the configuration file aka sort all the Targets.Range contents to correctly
	// represent a range. [0] > [1]
	for _, target := range m.Targets {
		if !(target.Range[0] < target.Range[1]) {
			sort.Ints(target.Range)
		}
	}

	// Return everything!
	return m, nil
}

package config

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"os"
)

// WARNING: never directly use this variable before calling getUserHome()
var configFilePath = "/etc/tfanc/tfanc.json"

// LoadConfig loads a configuration file from the standard path, defined by "configFilePath"
func LoadConfig() (Configuration, error) {

	f, err := os.Open(configFilePath)
	defer f.Close()

	if err != nil {
		return Configuration{}, err
	}

	fReader := bufio.NewReader(f)

	buf := new(bytes.Buffer)
	buf.ReadFrom(fReader)

	var m Configuration
	json.Unmarshal(buf.Bytes(), &m)

	//
	// Configuration file sanity checks
	//

	// The first condition is true only if "targets" in the json file is not even defined, so
	// we'll check for array length too.
	if m.Targets == nil || len(m.Targets) == 0 {
		return Configuration{}, errors.New("configuration error: no \"targets\" field defined, cannot continue without. Check your configuration file, it's malformed")
	}

	// Return everything!
	return m, nil
}

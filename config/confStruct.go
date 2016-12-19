package config

import "errors"

// Configuration is a structure holding content of the configuration file
type Configuration struct {
	Targets []Target `json:"targets"`
}

// Target is a structure holding a range and a level
type Target struct {
	MinTemp int `json:"mintemp"`
	Level   int `json:"level"`
}

// NextTarget returns the next useful target to the one passed as argument
func (c *Configuration) NextTarget(t Target) Target {
	for _, e := range c.Targets {
		if e.Level > t.Level {
			return e
		}
	}
	return t
}

// PrevTarget returns the previous useful target to the one passed as argument
func (c *Configuration) PrevTarget(t Target) Target {
	for _, e := range c.Targets {
		if e.Level == t.Level-1 {
			return e
		}
	}
	return t
}

// WhatLevelDoIBelong detects to what level a temp belongs.
func (c *Configuration) WhatLevelDoIBelong(temp int) (Target, error) {
	for index, element := range c.Targets {
		if temp < element.MinTemp {
			return c.Targets[index-1], nil
		}
	}

	if temp > c.Targets[len(c.Targets)-1].MinTemp {
		return c.Targets[len(c.Targets)-1], nil
	} else if temp < c.Targets[0].MinTemp {
		return c.Targets[0], nil
	}

	return Target{Level: 0, MinTemp: 0}, errors.New("cannot find a target for given temp")
}

type ByRange []Target

func (a ByRange) Len() int { return len(a) }

func (a ByRange) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByRange) Less(i, j int) bool {
	return a[i].MinTemp < a[j].MinTemp
}

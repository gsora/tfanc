package config

// Configuration is a structure holding content of the configuration file
type Configuration struct {
	Targets []Target `json:"targets"`
}

// Targets is a structure holding a range and a level
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

type ByRange []Target

func (a ByRange) Len() int { return len(a) }

func (a ByRange) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByRange) Less(i, j int) bool {
	return a[i].MinTemp < a[j].MinTemp
}

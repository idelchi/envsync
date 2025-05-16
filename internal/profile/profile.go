package profile

// Profile holds metadata plus an env-var map.
type Profile struct {
	Env     Env      `toml:"env,omitempty"     yaml:"env,omitempty"`
	DotEnv  []string `toml:"dotenv,omitempty"  yaml:"dotenv,omitempty"`
	Extends []string `toml:"extends,omitempty" yaml:"extends,omitempty"`
}

// newProfile creates a new profile with an empty env-var map.
func newProfile() *Profile {
	return &Profile{
		DotEnv:  []string{},
		Env:     make(Env),
		Extends: []string{},
	}
}

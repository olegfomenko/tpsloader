package config

type Config struct {
	AdminSeed  string            `yaml:"admin"`
	Passphrase string            `yaml:"passphrase"`
	HorizonURL string            `yaml:"horizon"`
	Creators   []string          `yaml:"creators,flow"`
	Payers     map[string]string `yaml:"payers,flow"`
	Amount     string            `yaml:"amount"`
	Duration   int64             `yaml:"duration"`
}

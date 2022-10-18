package definitions

type Config struct {
	Backup Backup   `yaml:"backup"`
	Policy []Policy `yaml:"policy"`
}

type Backup struct {
	Naming string `yaml:"naming"`
}

type Policy struct {
	From     string   `yaml:"from"`
	To       string   `yaml:"to"`
	Strategy Strategy `yaml:"strategy"`
}

type Strategy struct {
	Name   string  `yaml:"name"`
	Window *string `yaml:"window"`
}

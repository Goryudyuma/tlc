package tlc

type Config struct {
	ConsumerKey    string `yaml:"ConsumerKey"`
	ConsumerSecret string `yaml:"ConsumerSecret"`
	SeedString     string `yaml:"SeedString"`
}

type List struct {
	Listname          string
	Owner_screen_name string
	Owner_id          int64
}

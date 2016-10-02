package tlc

type MyTwitterKey struct {
	ConsumerKey       string `yaml:"ConsumerKey"`
	ConsumerSecret    string `yaml:"ConsumerSecret"`
	AccessToken       string `yaml:"AccessToken"`
	AccessTokenSecret string `yaml:"AccessTokenSecret"`
}

type List struct {
	Listname          string
	Owner_screen_name string
	Owner_id          int64
}

package tlc

type Config struct {
	ConsumerKey    string `yaml:"ConsumerKey"`
	ConsumerSecret string `yaml:"ConsumerSecret"`
	SeedString     string `yaml:"SeedString"`
}

type List struct {
	Listname          string `json:"Listname"`
	Owner_screen_name string `json:"OwnerScreenName"`
	Owner_id          int64  `json:"OwnerId"`
}

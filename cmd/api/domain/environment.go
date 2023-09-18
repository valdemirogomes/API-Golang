package domain

type MySQL struct {
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	Host          string `yaml:"host"`
	Database      string `yaml:"database"`
	Drive         string `yaml:"drive"`
	PoolSizeMax   int    `yaml:"poolsizemax"`
	PoolSizeIddle int    `yaml:"poolsizeiddle"`
}

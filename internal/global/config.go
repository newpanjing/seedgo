package global

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Dsn string `mapstructure:"dsn"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int64  `mapstructure:"expire"`
}

type AuthConfig struct {
	PublicPaths []string `mapstructure:"public_paths"`
}

type PermissionConfig struct {
	ExcludePaths []string `mapstructure:"exclude_paths"`
}

type Configuration struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Auth       AuthConfig       `mapstructure:"auth"`
	Permission PermissionConfig `mapstructure:"permission"`
}

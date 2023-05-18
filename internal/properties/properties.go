package properties

type SecurityProperties struct {
	Properties ServiceProperties `yaml:"security"`
}

type ServiceProperties struct {
	Grpc     GrpcProperties     `yaml:"grpc"`
	Cache    CacheProperties    `yaml:"cache"`
	Postgres PostgresProperties `yaml:"postgres"`
	Sso      SsoProperties      `yaml:"sso"`
}

type GrpcProperties struct {
	Network string `yaml:"network"`
	Port    string `yaml:"port"`
}

type CacheProperties struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

type PostgresProperties struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type SsoProperties struct {
	IntrospectionUrl string `yaml:"introspection_url"`
	TokenUrl         string `yaml:"token_url"`
	UserProfileUrl   string `yaml:"user_profile_url"`
	ClientId         string `yaml:"client_id"`
	ClientSecret     string `yaml:"client_secret"`
	Host             string `yaml:"host"`
}

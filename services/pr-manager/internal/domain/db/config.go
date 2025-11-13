package db

type Config struct {
	User     string
	Db       string
	Host     string
	Password string
}

func NewConfig(user string, db string, host string, password string) *Config {
	return &Config{User: user, Db: db, Host: host, Password: password}
}

func (c Config) GetConnectionString() string {
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":5432/" + c.Db
}

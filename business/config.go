package business

type WorkerConfig struct {
	BufferSize   int `env:"BUFFER_SIZE" envDefault:"1000"`
	FlushTimeout int `env:"FLUSH_TIMEOUT_MS" envDefault:"1000"`
}

type MailConfig struct {
	Host     string `env:"MAIL_SERVER_HOST"`
	Port     int    `env:"MAIL_SERVER_PORT"`
	Account  string `env:"MAIL_ACCOUNT"`
	Password string `env:"MAIL_PASSWORD"`
}

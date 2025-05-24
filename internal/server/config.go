package server

type ServerConfig struct {
	logLevel int
}

func NewServerConfig(logLevel int) *ServerConfig {
	return &ServerConfig{
		logLevel: logLevel,
	}
}

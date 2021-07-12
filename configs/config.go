package configs

import (
	"fmt"
)

// Config ...
type Config struct {
	HTTP        ConnAddress
	Tracer      ConnAddress
	ServiceName string
	Mongo       Mongo
	Storage     Storage
}

// ConnAddress ...
type ConnAddress struct {
	Host string
	Port string
}

func (c *ConnAddress) String() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// Mongo ...
type Mongo struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

// Storage ...
type Storage struct {
	AccessKey  string
	SecretKey  string
	BucketName string
	Region     string
}

func (m *Mongo) GetConnectString(host, port string) string {
	conn := host + ":" + port
	return conn
}

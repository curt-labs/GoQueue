package rabbitmq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Config struct {
	Username string
	Password string
	Hostname string
	Port     int
}

type ConsumerConfig struct {
	ExchangeName string `json:"exchange"`
	RoutingKey   string `json:"routing_key"`
	QueueName    string `json:"queue_name"`
	GATrackingID string `json:"ga_tracking"`
}

func NewConfig() *Config {
	c := new(Config)
	c.Hostname = "localhost"
	c.Port = 5672
	if os.Getenv("AMQP_HOST") != "" {
		c.Hostname = os.Getenv("AMQP_HOST")
	}
	if os.Getenv("AMQP_PORT") != "" {
		if port, err := strconv.Atoi(os.Getenv("AMQP_PORT")); err == nil {
			c.Port = port
		}
	}
	if os.Getenv("AMQP_USER") != "" {
		c.Username = os.Getenv("AMQP_USER")
	}
	if os.Getenv("AMQP_PASSWORD") != "" {
		c.Password = os.Getenv("AMQP_PASSWORD")
	}
	return c
}

func (c *Config) GetConnectionString() string {
	if c.Username != "" && c.Password != "" {
		return fmt.Sprintf("amqp://%s:%s@%s:%d", c.Username, c.Password, c.Hostname, c.Port)
	}
	return fmt.Sprintf("amqp://%s:%d", c.Hostname, c.Port)
}

func LoadConsumersConfig(filename string) (configs []ConsumerConfig, err error) {
	if filename == "" {
		err = fmt.Errorf("Error: %s", "Missing Configuration file!")
		return
	}
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(contents, &configs)

	return
}

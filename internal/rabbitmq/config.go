package rabbitmq

type Config struct {
	URI          string
	Exchange     string
	ExchangeType string
	Queue        string
	RoutingKey   string
}

module local/stocksbot

go 1.18

require (
	github.com/regismelgaco/go-sdks/logger v0.0.0-20221130114229-1ef5f784373a
	go.uber.org/zap v1.23.0
)

require (
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/rabbitmq/amqp091-go v1.5.0 // indirect
	github.com/regismelgaco/go-sdks/erring v0.0.0-20221127113222-947ccd31e2bf // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
)

replace github.com/regismelgaco/go-sdks/erring => ../../go-sdks/erring

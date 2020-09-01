module github.com/tOnkowzl/libs/saramax

go 1.15

replace (
	github.com/tOnkowzl/libs/contextx => ../contextx
	github.com/tOnkowzl/libs/logx => ../logx
)

require (
	github.com/Shopify/sarama v1.27.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/tOnkowzl/libs/logx v0.0.0-00010101000000-000000000000
)

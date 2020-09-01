module github.com/tOnkowzl/libs/pubsubx

go 1.15

replace (
	github.com/tOnkowzl/libs/contextx => ../contextx
	github.com/tOnkowzl/libs/logx => ../logx
)

require (
	cloud.google.com/go/pubsub v1.6.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/tOnkowzl/libs/logx v0.0.0-00010101000000-000000000000
	google.golang.org/api v0.30.0
)

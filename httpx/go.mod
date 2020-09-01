module github.com/tOnkowzl/libs/httpx

go 1.15

replace (
	github.com/tOnkowzl/libs/contextx => ../contextx
	github.com/tOnkowzl/libs/logx => ../logx
)

require (
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.6.1
	github.com/tOnkowzl/libs/contextx v0.0.0-00010101000000-000000000000
	github.com/tOnkowzl/libs/logx v0.0.0-00010101000000-000000000000
)

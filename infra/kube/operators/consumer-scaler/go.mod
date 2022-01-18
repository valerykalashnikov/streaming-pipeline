module github.com/valerykalashnikov/streaming-pipeline/infra/kube/operators/consumer-scaler

go 1.16

require (
	github.com/adjust/gocheck v0.0.0-20131111155431-fbc315b36e0e // indirect
	github.com/adjust/rmq v1.0.0
	github.com/adjust/uniuri v0.0.0-20130923163420-498743145e60 // indirect
	github.com/garyburd/redigo v1.6.3 // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	gopkg.in/bsm/ratelimit.v1 v1.0.0-20160220154919-db14e161995a // indirect
	gopkg.in/redis.v3 v3.6.4 // indirect
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
	sigs.k8s.io/controller-runtime v0.10.0
)

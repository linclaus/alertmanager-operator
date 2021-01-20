module github.com/linclaus/alertmanager-operator

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/prometheus/alertmanager v0.21.0
	github.com/prometheus/common v0.10.0
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	sigs.k8s.io/controller-runtime v0.5.0
)

module github.com/tektoncd/plumbing/catlin

go 1.14

require (
	github.com/google/go-containerregistry v0.4.1-0.20210128200529-19c2b639fab1
	github.com/spf13/cobra v1.0.0
	github.com/tektoncd/pipeline v0.26.0
	go.uber.org/zap v1.16.0
	gotest.tools v2.2.0+incompatible
	gotest.tools/v3 v3.0.2
	k8s.io/apimachinery v0.20.7
	k8s.io/client-go v0.20.7
	knative.dev/pkg v0.0.0-20210510175900-4564797bf3b7
)

// Pin k8s deps to 1.17.6
// replace (
// 	k8s.io/api => k8s.io/api v0.17.6
// 	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.6
// 	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
// 	k8s.io/apiserver => k8s.io/apiserver v0.17.6
// 	k8s.io/client-go => k8s.io/client-go v0.17.6
// 	k8s.io/code-generator => k8s.io/code-generator v0.17.6
// 	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20200410145947-bcb3869e6f29
// )

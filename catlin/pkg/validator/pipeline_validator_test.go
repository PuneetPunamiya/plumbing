package validator

import (
	"strings"
	"testing"

	"github.com/tektoncd/plumbing/catlin/pkg/parser"
	"gotest.tools/v3/assert"
)

const pipelineWithValidImageRef = `
apiVersion: tekton.dev/v1beta1
kind: Pipeline
metadata:
  name: foo
  labels:
    app.kubernetes.io/version: "0.1"
  annotations:
    tekton.dev/pipelines.minVersion: "0.12.1"
    tekton.dev/tags: bar1, bar2
    tekton.dev/displayName: "Foo"
    tekton.dev/platforms: linux/amd64
spec:
  description: >
    It is a foo pipeline
  serviceAccountName: pipeline
  tasks:
    - name: bar1
      taskRef:
        name: hello-world
        bundle: docker.io/foo/bar:0.1@sha256:053a6cb9f3711d4527dd0d37ac610e8727ec0288a898d5dfbd79b25bcaa29828
    - name: bar2
      taskSpec:
        params:
          - name: commit
        steps:
          - name: bar2
            image: docker.io/alpine:0.1@sha256:053a6cb9f3711d4527dd0d37ac610e8727ec0288a898d5dfbd79b25bcaa29828
            script: |
              if [[ ! $(params.commit) ]]; then
                exit 1
              fi
`

func TestPipelineValidator_ValidImageRef(t *testing.T) {

	r := strings.NewReader(pipelineWithValidImageRef)
	parser := parser.ForReader(r)

	res, err := parser.Parse()
	assert.NilError(t, err)

	v := ForKind(res)
	result := v.Validate()
	assert.Equal(t, 0, result.Errors)

	lints := result.Lints
	assert.Equal(t, 0, len(lints))
}

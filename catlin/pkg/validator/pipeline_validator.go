package validator

import (
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/tektoncd/plumbing/catlin/pkg/parser"
)

type pipelineValidator struct {
	res *parser.Resource
}

var _ Validator = (*pipelineValidator)(nil)

func NewPipelineValidator(r *parser.Resource) *pipelineValidator {
	return &pipelineValidator{res: r}
}

func (p *pipelineValidator) Validate() Result {

	result := Result{}

	res, err := p.res.ToType()
	if err != nil {
		result.Error("Failed to decode to a Pipeline - %s", err)
		return result
	}

	pipeline := res.(*v1beta1.Pipeline)
	for _, task := range pipeline.Spec.Tasks {
		if task.TaskSpec != nil {
			for _, step := range task.TaskSpec.Steps {
				result.Append(validateStep(step))
			}
		}
	}
	for _, task := range pipeline.Spec.Tasks {
		if task.TaskRef != nil {
			if task.TaskRef.Bundle != "" {
				result.Append(p.validateBundle(task.TaskRef.Bundle))
			}
		}
	}
	for _, task := range pipeline.Spec.Finally {
		if task.TaskSpec != nil {
			for _, step := range task.TaskSpec.Steps {
				result.Append(validateStep(step))
			}
		}
	}

	return result
}

func (p *pipelineValidator) validateBundle(bundle string) Result {
	result := Result{}

	if _, usesVars := extractExpressionFromString(bundle, ""); usesVars {
		result.Warn("Task uses image %q that contains variables; skipping validation", bundle)
		return result
	}

	if !strings.Contains(bundle, "/") || !isValidRegistry(bundle) {
		result.Warn("Task uses image %q; consider using a fully qualified name - e.g. docker.io/library/ubuntu:1.0", bundle)
	}

	if strings.Contains(bundle, "@sha256") {
		rep, err := name.NewDigest(bundle, name.WeakValidation)
		if err != nil {
			result.Error("Task uses image %q with an invalid digest. Error: %s", bundle, err)
			return result
		}

		if !tagWithDigest(rep.String()) {
			result.Warn("Task uses image %q; consider using a image tagged with specific version along with digest eg. abc.io/img:v1@sha256:abcde", bundle)
		}
		return result
	}

	ref, err := name.NewTag(bundle, name.WeakValidation)
	if err != nil {
		result.Error("Task uses image %q with an invalid tag. Error: %s", bundle, err)
		return result
	}

	if strings.Contains(ref.Identifier(), "latest") {
		result.Error("Task uses image %q which must be tagged with a specific version", bundle)
	}

	return result
}

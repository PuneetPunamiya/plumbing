package create

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/tektoncd/plumbing/catlin/pkg/app"
	"github.com/tektoncd/plumbing/catlin/pkg/parser"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

type options struct {
	cli     app.CLI
	from    string
	version string
	kind    string
	args    []string
}

func Command(p app.CLI) *cobra.Command {
	opts := &options{cli: p}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create resources",
		Long:    "Creates a new resource in the catalog",
		Example: "catlin create",
	}

	cmd.AddCommand(
		commandForKind("task", opts),
	)

	return cmd
}

// commandForKind creates a cobra.Command that when run sets
// opts.Kind and opts.Args and invokes opts.run
func commandForKind(kind string, opts *options) *cobra.Command {

	cmd := &cobra.Command{
		Use:          kind,
		Short:        "Kind of resource is " + strings.Title(kind),
		Long:         ``,
		SilenceUsage: true,
		Example:      "catlin create task",
		RunE:         opts.createResource,
	}

	return cmd
}

func (opts *options) createResource(cmd *cobra.Command, args []string) error {

	if len(args) < 2 {
		log.Fatal("Please provide a path to create a resource")
	}

	// Get the path to the directory
	path := args[1]
	_, err := os.Stat(path)
	if err == nil {
		log.Fatal("Resource exists!! Please bump the version of the resource using `catlin bump` command")
	}
	if os.IsNotExist(err) {

		if err := os.MkdirAll(args[1], os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	// Open and read the template yaml
	r, err := os.Open("pkg/cmd/create/template.yaml")
	if err != nil {
		return err
	}

	// Parse the template manifest
	parser := parser.ForReader(r)
	res, err := parser.Parse()
	if err != nil {
		return err
	}

	// Get the resource
	t := &v1beta1.Task{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(res.Unstructured.Object, t); err != nil {
		return err
	}

	// Marshal the resource to save the file in the specified location
	yaml, err := yaml.Marshal(t)

	err = ioutil.WriteFile(filepath.Join(args[1], args[0]+".yaml"), yaml, 0666)
	if err != nil {
		return err
	}

	return nil
}

package test_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/redhat-et/copilot-ops/pkg/cmd"
)

func TestGeneratePodForPVC(t *testing.T) {
	t.Log(os.Getwd())
	Run(t, []string{
		"generate",
		"--file",
		"examples/app1/mysql-pvc.yaml",
		"--request",
		"Generate a pod that mounts the PVC. Set the pod resources requests and limits to 4 cpus and 5 Gig of memory.",
	})
}

// TestEditPVCSize Tests that the edit command can successfully change the size of a PVC YAML
// to be 100Gi.
func TestEditPVCSize(t *testing.T) {
	Run(t, []string{
		"edit",
		"--file",
		"examples/app1/mysql-pvc.yaml",
		"--request",
		"Increase the size of the PVC to 100Gi.",
	})
}

func Run(t *testing.T, args []string) string {
	cmd := cmd.NewRootCmd()
	buf := bytes.NewBufferString("")
	cmd.SetOut(buf)
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}
	out := string(bytes)
	t.Logf("out: %+v\n", out)
	return out
}

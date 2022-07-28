package cmd_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-et/copilot-ops/pkg/openai"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

// OpenAITestServer Creates a mocked OpenAI server which can pretend to handle requests during testing.
func OpenAITestServer() *httptest.Server {
	return httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resBytes []byte
		switch r.URL.Path {
		case "/v1/edits":
			res := openai.EditResponse{
				Response: openai.Response{
					Choices: []openai.Choice{
						{
							Text: `
# @path/to/kubernetes.yaml
apiVersion: v1
kind: Pod
metadata:
	name: cute-cats
spec:
	priority: high
`,
							Index: 0,
						},
					},
				},
			}
			resBytes, _ = json.Marshal(res)
			fmt.Fprint(w, string(resBytes))
			return
		case "/v1/completions":
			res := openai.CompletionResponse{
				ID: "completion id",
				Response: openai.Response{
					Choices: []openai.Choice{
						{
							Text:  "choice 1",
							Index: 0,
						},
					},
				},
			}
			resBytes, _ = json.Marshal(res)
			fmt.Fprintln(w, string(resBytes))
			return
		default:
			// the endpoint doesn't exist
			http.Error(w, "the resource path doesn't exist", http.StatusNotFound)
			return
		}
	}))
}

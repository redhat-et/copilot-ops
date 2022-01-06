/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	completev1alpha1 "github.com/redhat-et/openshift-copilot-poc/api/v1alpha1"
)

var OPENAI_API_KEY string

// CompletionReconciler reconciles a Completion object
type CompletionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// create json body
type OpenAIRequestBody struct {
	Prompt           string   `json:"prompt"`
	Temperature      float64  `json:"temperature"`
	MaxTokens        int      `json:"max_tokens"`
	TopP             float64  `json:"top_p"`
	FrequencyPenalty float64  `json:"frequency_penalty"`
	PresencePenalty  float64  `json:"presence_penalty"`
	Stop             []string `json:"stop"`
}

type OpenAIResponseBody struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
}

//+kubebuilder:rbac:groups=complete.yaml-copilot.com,resources=completions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=complete.yaml-copilot.com,resources=completions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=complete.yaml-copilot.com,resources=completions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Completion object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *CompletionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// update the Completion status
	obj := &completev1alpha1.Completion{}

	// get the Completion object
	if err := r.Client.Get(ctx, req.NamespacedName, obj); err != nil {
		fmt.Printf("Error getting Completion instance: %v\n", err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// check if the completion object doesn't need to be reconciled
	if !obj.NeedsReconcile() {
		return ctrl.Result{}, nil
	}

	completion, err := createCompletion(obj.Spec.UserPrompt)
	if err != nil {
		return ctrl.Result{}, err
	}

	// request openai api here
	obj.Status.Completion = completion
	obj.Status.ObservedGeneration = obj.GetObjectMeta().GetGeneration()

	// update the status of the custom resource
	// err = r.Status().Update(ctx, completion)
	statusErr := r.Client.Status().Update(ctx, obj)
	if err == nil { // Don't mask previous error
		err = statusErr
	}

	return ctrl.Result{}, err
}

// get completion from openai
func createCompletion(prompt string) (string, error) {
	/*
			Create an HTTPS request from Golang based on the following cURL command:

			curl https://api.openai.com/v1/engines/davinci-codex/completions \
		  -H "Content-Type: application/json" \
		  -H "Authorization: Bearer $OPENAI_API_KEY" \
		  -d '{
		  "prompt": "",
		  "temperature": 0,
		  "max_tokens": 64,
		  "top_p": 1,
		  "frequency_penalty": 0,
		  "presence_penalty": 0,
		  "stop": ["'#'", "\\n\\n"]
		}'

	*/
	const url = "https://api.openai.com/v1/engines/davinci-codex/completions"
	var OPENAI_API_KEY = os.Getenv("OPENAI_API_KEY")

	var requestBody = OpenAIRequestBody{
		Prompt:           prompt,
		Temperature:      0.1,
		MaxTokens:        1024,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		Stop:             []string{"\n\n"},
	}

	// marshal above json body into a string
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	// tostring the json body
	body := io.Reader(bytes.NewReader(jsonBody))
	req, err := http.NewRequest("POST", url, body)

	// set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+OPENAI_API_KEY)

	// make an HTTPS POST request to the OpenAI API completions endpoint
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	// read the body with answer
	var answer OpenAIResponseBody
	data, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(data, &answer)
	if err != nil {
		fmt.Println(err)
	}

	// return the answer
	return answer.Choices[0].Text, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CompletionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&completev1alpha1.Completion{}).
		Complete(r)
}

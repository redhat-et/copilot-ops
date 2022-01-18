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
	"reflect"
	"strings"

	// "gopkg.in/yaml.v2"

	"gopkg.in/yaml.v2"
	k8syaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	// appsv1 "k8s.io/api/apps/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
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
	logger := log.FromContext(ctx)

	logger.Info("Getting kubeconfig & starting dynamic client")
	kubeConfig := ctrl.GetConfigOrDie()
	dynamicClient, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		logger.Info("Error creating dynamic client", "error", err)
	}

	// update the Completion status
	obj := &completev1alpha1.Completion{}

	// get the Completion object
	if err := r.Client.Get(ctx, req.NamespacedName, obj); err != nil {
		logger.Error(err, "Error getting Completion instance")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// check if the completion object doesn't need to be reconciled
	if !obj.NeedsReconcile() {
		return ctrl.Result{}, nil
	}

	completion, err := createCompletion(obj)
	if err != nil {
		return ctrl.Result{}, err
	}
	logger.Info("Completion created", "completion", completion)
	obj.Status.Completions = strings.Split(completion, "---")
	obj.Status.ObservedGeneration = obj.GetObjectMeta().GetGeneration()

	// create a list to hold the maps for the unmarshaled yamls
	// var yamlList []map[string]interface{}

	// split the obtained completion & try to decode it into YAML objects so we can try to create them through the client
	completionSplit := strings.Split(completion, "---")
	for _, generatedYaml := range completionSplit {
		logger.Info("now processing generated YAML", "generatedYaml", generatedYaml)
		if generatedYaml == "" {
			continue
		}
		// deconstruct the string into yaml and attempt to create it through the k8s api
		m := make(map[interface{}]interface{})

		err := yaml.Unmarshal([]byte(generatedYaml), &m)
		if err != nil {
			logger.Error(err, "could not unmarshal the yaml")
			continue
		}
		logger.Info("unmarshalled yaml", "yaml", m)
		// obtain the group and kind from the yaml
		apiGroup, apiVersion := getApiGroupAndVersion(m["apiVersion"].(string))
		kind := m["kind"].(string)
		logger.Info("attempting to create object of the types", "apiGroup", apiGroup, "apiVersion", apiVersion, "kind", kind)
		gvr := schema.GroupVersionResource{
			Group:    apiGroup,
			Version:  apiVersion,
			Resource: strings.ToLower(kind) + "s",
		}

		// convert from singular to plural
		logger.Info("Creating object", "gvr", gvr)
		// create the object through the client
		decUnstructured := k8syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
		unstructuredK8sResource := &unstructured.Unstructured{}
		_, _, err = decUnstructured.Decode([]byte(generatedYaml), nil, unstructuredK8sResource)
		if err != nil {
			logger.Error(err, "Error decoding generated YAML")
		}

		_, err = dynamicClient.Resource(gvr).Namespace("default").Create(context.Background(), unstructuredK8sResource, v1.CreateOptions{})
		if err != nil {
			logger.Error(err, "cannot create the resource")
		}

		// if generatedYaml == "" {
		// 	continue
		// }
		// // deconstruct the string into yaml and attempt to create it through the k8s api
		// m := make(map[interface{}]interface{})

		// err := yaml.Unmarshal([]byte(generatedYaml), &m)
		// if err != nil {
		// 	logger.Error(err, "could not unmarshal the yaml")
		// 	continue
		// }

		// // convert m to be of type map[string]interface{}
		// nm := make(map[string]interface{})
		// for k, v := range m {
		// 	nm[k.(string)] = v
		// }

		// // add to list
		// yamlList = append(yamlList, nm)
	}

	// try and apply the YAMLs
	// ApplyYamls(ctx, r.Client, dynamicClient, yamlList)
	statusErr := r.Client.Status().Update(ctx, obj)
	if err == nil { // Don't mask previous error
		err = statusErr
	}

	return ctrl.Result{}, err
}

// get completion from openai
func createCompletion(obj *completev1alpha1.Completion) (string, error) {
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

	var fullUserPrompt = `
## Below is a comment describing a desired deployment in Kubernetes, followed by YAMLs to create the described setup

# Description: ` + obj.Spec.UserPrompt + "\n---"

	// fmt.Printf("fullUserPrompt:---\n%s\n---", fullUserPrompt)

	var requestBody = OpenAIRequestBody{
		Prompt:           fullUserPrompt,
		Temperature:      0.1,
		MaxTokens:        obj.Spec.MaxTokens,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		// Stop:             []string{""},
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

func ApplyYamls(ctx context.Context, client client.Client, dc dynamic.Interface, yamlList []map[string]interface{}) error {
	// gvr := schema.GroupVersionResource{
	// 	Group:    "",
	// 	Version:  "v1",
	// 	Resource: "pods",
	// }
	// pods, err := dc.Resource(gvr).Namespace("kube-system").List(context.Background(), v1.ListOptions{})
	// if err != nil {
	// 	fmt.Printf("error getting pods: %v\n", err)
	// 	os.Exit(1)
	// }

	// for _, pod := range pods.Items {
	// 	fmt.Printf(
	// 		"Name: %s\n",
	// 		pod.Object["metadata"].(map[string]interface{})["name"],
	// 	)
	// }

	// go through our list of yaml resources & attempt to apply each one
	for _, rsrc := range yamlList {
		// list of failpoints
		var failpoints map[string]bool = map[string]bool{}

		kind, ok := rsrc["kind"]
		failpoints["kind"] = ok

		apiVersionStr, ok := rsrc["apiVersion"]
		failpoints["apiVersion"] = ok

		spec, ok := rsrc["spec"]
		failpoints["spec"] = ok

		metadata, ok := rsrc["metadata"]
		failpoints["metadata"] = ok

		// make sure metadata is not empty
		if metadata == nil {
			metadata = map[string]interface{}{}
		}

		// create the yaml object
		obj := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       kind,
				"apiVersion": apiVersionStr,
				"spec":       spec,
				"metadata":   metadata,
			},
		}

		fmt.Printf("current values for object: %+v\n", obj)
		// get types of all fields in the map
		for k, v := range rsrc {
			fmt.Printf("key: %s, value: %+v\n", k, reflect.TypeOf(v))
		}

		// set labels
		_, ok = obj.Object["metadata"]
		if !ok || obj.Object["metadata"] == nil {
			obj.Object["metadata"] = map[interface{}]interface{}{}
		}
		labels, ok := obj.Object["metadata"].(map[interface{}]interface{})["labels"]
		fmt.Printf("labels: %+v\n", labels)
		if !ok {
			obj.Object["metadata"].(map[interface{}]interface{})["labels"] = map[interface{}]interface{}{}
		}
		obj.Object["metadata"].(map[interface{}]interface{})["labels"].(map[interface{}]interface{})["generatedBy"] = "copilot-operator"

		// set GVR
		gvr := schema.GroupVersionResource{
			Group:    obj.GroupVersionKind().Group,
			Version:  obj.GroupVersionKind().Version,
			Resource: obj.GroupVersionKind().Kind,
		}
		// decUnstructured := k8syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

		// testMachineConfigPools := &unstructured.Unstructured{}

		// _, _, err := decUnstructured.Decode([]byte(machineconfigpoolYAML), nil, testMachineConfigPools)

		obj, err := dc.Resource(gvr).Create(ctx, obj, v1.CreateOptions{})
		if err != nil {
			fmt.Printf("error creating object: %v\n", err)
			continue
		} else {
			fmt.Printf("created object: %v\n", obj)
		}

		// fmt.Printf("resourceId: %+v\n", resourceId)
	}

	return nil
}

func getApiGroupAndVersion(apiVersionStr string) (string, string) {
	// split apiversion into group and version
	var apiVersion, apiGroup string

	apiVersionSplit := strings.Split(apiVersionStr, "/")
	if len(apiVersionSplit) != 2 {
		fmt.Printf("apiVersion %s is not in the format group/version, this must be a core object", apiVersionStr)
		apiGroup = ""
		apiVersion = apiVersionSplit[0]
	} else {
		apiGroup = apiVersionSplit[0]
		apiVersion = apiVersionSplit[1]
	}
	return apiGroup, apiVersion
}

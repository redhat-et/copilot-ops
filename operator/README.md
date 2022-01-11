# OpenShift Copilot Operator

This is a proof-of-concept Kubernetes operator which
provides completions to the user through the OpenAI Codex API. 

## About

The `Completion` CRD is used to interface with the operator, and
the `.spec.userPrompt` is used to provide a prompt to the operator which
the operator will then attempt to complete.

E.g.:

The following Completion resource defines a prompt to be completed:

```yaml
apiVersion: copilot.poc.com/v1
kind: Completion
metadata:
  name: completion
spec:
  userPrompt: "# Deploy a Django server with a MySQL deployment which stores data in a 20Gi PVC\nkind: Deployment\napiVersion:"

```

And the operator will return the result in `.status.completion`:

```
apiVersion: copilot.poc.com/v1
kind: Completion
metadata:
  name: completion
spec:
  userPrompt: "# Deploy a Django server with a MySQL deployment which stores data in a 20Gi PVC\nkind: Deployment\napiVersion:"
status:
	completion: |
		apiVersion: apps/v1
		kind: Deployment
		metadata:
			name: django-deployment
		spec:
			replicas: 1
				selector:
					matchLabels:
						app: django
				template:
					metadata:
						labels:
							app: django
						spec:
							containers:
							- name: django
								image: django-image
								ports:
			# ...
```

## Installation

To install this operator, you will first need to have Golang version 16 installed.

Then you can install the operator with:

```sh
make generate
make manifests
```

## Usage

To start the operator, simply run `make install run`


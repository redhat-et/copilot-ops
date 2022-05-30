# copilot-ops

`copilot-ops` is a CLI tool that boosts up any "devops repo" to a ninja level of *Artificially Intelligent Ops Repo*.

## Requirements

In order to use `copilot-ops`, you need to have an OpenAI account with access to the GPT-3 Codex model,
and an API token saved as the `OPENAI_API_TOKEN` environment variable.

## Installation

Installing `copilot-ops` is simple, simply clone this repository and `make build`:

```console
# Clone the repository
git clone https://github.com/redhat-et/copilot-ops.git

# Build the binary
make build
```



## Usage 

`copilot-ops` currently supports two functionalities: `generate` and `edit`. 

By default, `copilot-ops` will only print to stdout. To write the
changes directly to the disk, provide the `--write` flag.


### Editing Files

`copilot-ops` is capable of updating existing files using a command phrased with natural language.
To do this, you would simply reference your files as described above, and then specify the necessary changes.

For example:

```bash
copilot-ops edit --request="The timeout for the @server.yaml should be increased from 60s to 2m" \
  --file deployments/server.yaml`
```


### Generating Files

The `generate` command accepts a description of the file(s) needed and a set of files which are used to generate a new file based on their contents.

Here are a few examples with the `generate` command:

```sh
# to create a Jupyter Notebook which uses a GPU to accelerate machine-learning tasks
copilot-ops generate --request "create a Deployment which pulls a Jupyter Notebook image and requests 1 GPU resource"

# to generate a Pod that mounts a given ConfigMap
copilot-ops generate -f examples/stock-data.yaml --request '
	create a Pod which runs an express.js app and mounts the stock-data ConfigMap to trade stocks
'

# launch a Job which pulls data from the S3 bucket at 's3://my-bucket/data.csv' and loads it into a PVC in the same namespace
copilot-ops generate -f examples/aws-credentials-secret.yaml --request '
	create a Job which pulls data from the S3 bucket at "s3://my-bucket/data.csv" and loads it into a PVC in the same namespace
'
```

To control the amount of tokens used when generating, you can also
specify the `--ntokens` flag.

Here's an example where we want OpenAI to generate a service based
on a Deployment, but should not exceed 100 tokens:

```bash
copilot-ops generate --request "Create a Service named 'mongodb-service' to expose the mongodb-deployment" \
	--file deployments/mongodb-deployment.yaml \
	--ntokens 100
```

To avoid providing multiple files, we can use the `--filesets` flag to specify a list of filesets to use.

For example:

```bash
copilot-ops generate --request "Create a Service for each of these deployments" --fileset deployments
```

### Under the hood

In a nutshell, `copilot-ops` functions by formatting the user input and provided files, if any, in a way that an OpenAI would understand it as a programmer taking an issue and updating it.

#### Generating Files

Here's a breakdown of the generate process:

1. Format the user's request along with the necessary files using a generate template that will be "autocompleted" by OpenAI.
1. Have OpenAI attempt to complete the prompt and retrieve OpenAI's response.
1. If OpenAI succeeded, parse the response and extract the newly generated files.
1. Write the generated files either to the disk or to STDOUT. 


#### Editing Files

Editing files is similar, however it currently only works on one file at a time. 

A breakdown of the process is described below:

1. Take an input text description in natural language for a change requested on the repo. There should be no significant restrictions on the text format, but we will provide guidelines in the form of issue templates and identify textual markers to provide higher success rates and consistency.
1. Collect information from the repo to be attached to the AI request.
1. Send the request to AI service to process it and reply with the changed files.
1. Apply the reply changes back to the repo.

```console
> copilot-ops edit "set mysql memory to 42 Gi only for production env" --file deployments/mysql.yaml --write

Collecting deployments/mysql.yaml ........ OK
Using the force .......................... OK
Applying changes ......................... OK

Done
```



## Copilot-Ops CLI

The Copilot-Ops CLI applies AI principles in order to turn you into a GitOps ninja.
Provide an issue of the changes you'd like to make, and optionally some files that 
you'd like to reference in your issue, and the CLI will generate the according changes for you.


## About copilot-ops

The idea is that using modern AI services (such as [OpenAI Codex](https://openai.com/blog/openai-codex/)) users can describe their requested changes in a github issue, which triggers the `copilot-ops` action, which will use the AI service to create a suggested PR to resolve the issue. Users can then interact with the PR, review, make code changes directly, approve and merge it.

The assumption is that these new AI capabilities to understand natural language (NLP) and modify code, can simplify working on a devops repo where many changes are eventually simply modifying kubernetes yamls, config files or deployment scripts.

Given a request in natural language and access to the cloned repo, the `copilot-ops` tool works by collecting information from the repo, and composing a well formatted request to the AI service that has all the information to prepare the code change. When the response is received, it applies the changes back to the repo by mapping the reply to source file changes.

Runtime consideration - to be able to integrate this tool to various gitops frameworks without dragging along dependencies or forcing to run inside a container image, we decided to package this functionality in a standalone self-contained golang CLI for easy portability. This decision is not critical for the usability of the tool and can be revised in the future as needed.


Currently, Copilot-Ops is using OpenAI codex as a backend, but we plan to modularize this so that you may use another backend such DeepMind's AlphaCode, or IBM Watson.


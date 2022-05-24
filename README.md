# copilot-ops

`copilot-ops` is a CLI tool that boosts up any "devops repo" to a ninja level of *Artificially Intelligent Ops Repo*.

The idea is that using modern AI services (such as [OpenAI Codex](https://openai.com/blog/openai-codex/)) users can describe their requested changes in a github issue, which triggers the `copilot-ops` action, which will use the AI service to create a suggested PR to resolve the issue. Users can then interact with the PR, review, make code changes directly, approve and merge it.

The assumption is that these new AI capabilities to understand natural language (NLP) and modify code, can simplify working on a devops repo where many changes are eventually simply modifying kubernetes yamls, config files or deployment scripts.

Given a request in natural language and access to the cloned repo, the `copilot-ops` tool works by collecting information from the repo, and composing a well formatted request to the AI service that has all the information to prepare the code change. When the response is received, it applies the changes back to the repo by mapping the reply to source file changes.

Runtime consideration - to be able to integrate this tool to various gitops frameworks without dragging along dependencies or forcing to run inside a container image, we decided to package this functionality in a standalone self-contained golang CLI for easy portability. This decision is not critical for the usability of the tool and can be revised in the future as needed.

## `suggest` command

The `suggest` command is the most basic functionality which works like this:

1. Takes an input text description in natural language for a change requested on the repo. There should be no significant restrictions on the text format, but we will provide guidelines in the form of issue templates and identify textual markers to provide higher success rates and consistency.
1. Collects information from the repo to be attached to the AI request.
1. Sends the request to AI service to process it and reply with the changed files.
1. Apply the reply changes back to the repo.

```console
> copilot-ops suggest "set mysql memory to 42 Gi only for production env"
Collecting 2 files ........ OK
Using the force ........... OK
Applying 2 changes ........ OK
Done ...................... use `git diff/add/commit/etc`
```

## `generate` command

The `generate` command accepts a description of the file(s) needed and a set of files which are used to generate a new file based on their contents.

Here's a breakdown of the generate process:
1. Format the user's request along with the necessary files using a generate template that will be "autocompleted" by OpenAI.
1. Have OpenAI attempt to complete the prompt and retrieve OpenAI's response.
1. If OpenAI succeeded, parse the response and extract the newly generated files.
1. Write the generated files either to the disk or to STDOUT. 

### Usage 

```sh
# to generate a Pod that mounts a given ConfigMap
copilot-ops 
```

```console
> copilot-ops generate -f /path/to/file1.yaml -f /path/to/file2.yaml -f /path/to/file3.yaml
```

## Copilot-Ops CLI

The Copilot-Ops CLI applies AI principles in order to turn you into a GitOps ninja.
Provide an issue of the changes you'd like to make, and optionally some files that 
you'd like to reference in your issue, and the CLI will generate the according changes for you.

### Requirements

In order to use `copilot-ops`, you need to have an OpenAI account with access to the GPT-3 Codex model,
and an API token saved as the `OPENAI_API_TOKEN` environment variable.


### Basic Usage

Basic example:
`copilot-ops suggest --issue="Create a Deployment for a NodeJS server exposing a GraphQL API, and expose it through a ClusterIP"`

This will cause copilot-ops to generate a set of Kubernetes resources described as YAML files which could be deployed directly into your cluster.
In the case of new files, `copilot-ops` will attempt to write them into a `generated-files` directory to your directory's basepath.
If `--dry-run` is specified, then `copilot-ops` simply prints the generated YAML files to STDOUT, and no changes are made.


#### Referencing Files

Copilot-Ops is fully realized when referencing existing YAMLs in the description of your issue.
To do this, simply import them using the `-f` command using the format `-f @filetag:path/to/file.yaml`, and then reference them in your issue using `@filetag`.

For example:
`copilot-ops suggest --issue="Create a service for @nodejs exposing the GraphQL API through ClusterIP" -f @nodejs:deployments/nodejs-server.yaml`


#### Updating Files

Beyond creating new files, Copilot-Ops is capable of processing your issue and using it to update existing files.
To do this, you would simply reference your files as described above, and then specify the necessary changes.

For example:
`copilot-ops suggest --issue="The timeout for the @nodejs server should be increased from 60s to 2m" -f @nodejs:deployments/nodejs-server.yaml`



Currently, Copilot-Ops is using OpenAI codex as a backend, but we plan to modularize this so that you may use another backend such DeepMind's AlphaCode, or IBM Watson.


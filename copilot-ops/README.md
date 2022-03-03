## Copilot-Ops CLI

The Copilot-Ops CLI applies AI principles in order to turn you into a GitOps ninja.
Provide an issue of the changes you'd like to make, and optionally some files that 
you'd like to reference in your code, and the CLI will generate these changes for you.

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


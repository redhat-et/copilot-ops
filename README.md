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

## `suggest-pr` command

Same as `suggest` but also interacts with the source issue and create a target PR:

```console
> copilot-ops suggest-pr --issue=123
Reading issue #123 ........ OK
Collecting 2 files ........ OK
Using the force ........... OK
Applying 2 changes ........ OK
Creating the PR ........... OK
Done ...................... PR #456
```

## `update-pr` command

Similar to `suggest-pr` but as a followup to a new user comment on a PR that contains a change request on existing PR code:

```console
> copilot-ops update-pr --pr=456 --comment=789
Reading PR #456 ........... OK
Collecting 2 files ........ OK
Using the force ........... OK
Applying 2 changes ........ OK
Updating the PR ........... OK
Done ...................... PR #456
```



import Operator, { ResourceEvent, ResourceEventType } from "@dot-i/k8s-operator";
import axios from "axios";
import { OPENAI_API_KEY } from "./constants";
import { CompletionResource } from "./types";

export default class CompletionOperator extends Operator {
  group = "copilot.poc.com";
  version = "v1";
  plural = "completions";

  constructor() {
    super(/* pass in optional logger*/);
  }

  protected async init() {
    // NOTE: we pass the plural name of the resource
    await this.watchResource(this.group, this.version, this.plural, async (e) => {
      if (e.type === ResourceEventType.Added || e.type === ResourceEventType.Modified) {
        if (
          !(await this.handleResourceFinalizer(e, "mycustomresources.dot-i.com", (event) =>
            this.resourceDeleted(event),
          ))
        ) {
          await this.resourceModified(e);
        }
      }
    }).catch((err) => {
      console.log(err);
    });
  }

  private async resourceModified(e: ResourceEvent) {
    const object = e.object as CompletionResource;
    // const metadata = object.metadata;
    const { metadata, status } = object;
    console.log("metadata", metadata);
    console.log("e.meta", e.meta);

    if (!status || status.observedGeneration !== metadata.generation) {
      console.log("modifying", metadata.name);
      console.log("object", object);
      const openAIUrl = "https://api.openai.com/v1/engines/davinci-codex/completions";
      // request openai API to complete data

      const headers = {
        "Authorization": `Bearer ${OPENAI_API_KEY}`,
        "Content-Type": "application/json",
      };
      // console.log("headers", headers);
      const body = {
        prompt:
          "# Below is a series of YAML files used to create resources in a Kubernetes cluster\n" +
          object.spec.userPrompt,
        max_tokens: 50,
        stop: ["#\n#\n", "\n\n---\n\n", "\n\n"],
        temperature: 0.12,
        top_p: 1,
        frequency_penalty: 0,
        presence_penalty: 0,
      };

      let completion;
      // request the openai api using axios
      await axios.post(openAIUrl, body, { headers }).then(async (response) => {
        // update the object with the competion result
        if (response.status == 200 && response.data.choices) {
          console.log("before sending, the object is ", object);
          completion = response.data.choices[0]!.text;
          // console.log("received completion from openai: ", completion);
          // console.log("now we update the objects completion", object);
        }
      });

      await this.setResourceStatus(e.meta, {
        observedGeneration: object.metadata.generation,
        completion: completion,
      }).catch((err) => {
        console.log(err);
      });
    }
  }

  private async resourceDeleted(e: ResourceEvent) {
    // TODO: handle resource deletion here
    console.log("handling resource deletion: ", e.object);
  }
}

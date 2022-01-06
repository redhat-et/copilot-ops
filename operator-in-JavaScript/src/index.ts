/* eslint-disable @typescript-eslint/no-explicit-any */
import * as k8s from "@kubernetes/client-node";
import * as fs from "fs";
import axios from "axios";

const COMPLETION_GROUP = "copilot.poc.com";
const COMPLETION_VERSION = "v1";
const COMPLETION_PLURAL = "completions";

// get OPENAI key
const OPENAI_API_KEY = process.env["OPENAI_API_KEY"];
console.log("openai key", OPENAI_API_KEY);

type CompletionSpec = {
  userPrompt: string;
};
type CompletionStatus = {
  completion: string;
};

type Completion = {
  apiVersion: string;
  kind: string;
  metadata: k8s.V1ObjectMeta;
  spec?: CompletionSpec;
  status?: CompletionStatus;
};

const kc = new k8s.KubeConfig();
kc.loadFromDefault();

const k8sApiMC = kc.makeApiClient(k8s.CustomObjectsApi);

const watch = new k8s.Watch(kc);

async function onEvent(phase: string, apiObj: any) {
  log(`Received event in phase ${phase}.`);
  if (phase == "ADDED") {
    // ensure the status field exists
    apiObj.status = {
      completion: "",
    };
    scheduleReconcile(apiObj);
  } else if (phase == "MODIFIED") {
    scheduleReconcile(apiObj);
  } else if (phase == "DELETED") {
    console.log(`Deleted ${apiObj.metadata.name}`);
  } else {
    log(`Unknown event type: ${phase}`);
  }
}

function onDone(err: any) {
  console.log(`Connection closed`);
  if (typeof err !== "undefined") {
    console.log("got err: ", err);
  }
  watchResource();
}

async function watchResource(): Promise<any> {
  log("Watching API");
  return watch.watch(`/apis/${COMPLETION_GROUP}/${COMPLETION_VERSION}/${COMPLETION_PLURAL}`, {}, onEvent, onDone);
}

let reconcileScheduled = false;

function scheduleReconcile(obj: Completion) {
  if (!reconcileScheduled) {
    setTimeout(reconcileNow, 1000, obj);
    reconcileScheduled = true;
  }
}

async function reconcileNow(obj: Completion) {
  reconcileScheduled = false;
  const userData = obj.spec?.userPrompt;
  if (typeof userData !== "string") {
    console.error("user data is not a string");
    return;
  }

  try {
    const openAIUrl = "https://api.openai.com/v1/engines/davinci-codex/completions";
    // request openai API to complete data

    const headers = {
      "Authorization": `Bearer ${OPENAI_API_KEY}`,
      "Content-Type": "application/json",
    };
    const body = {
      prompt: "# Below is a series of YAML files used to create resources in a Kubernetes cluster\n" + userData,
      max_tokens: 1024,
      stop: ["#\n#\n", "\n\n---\n\n", "\n\n"],
      temperature: 0.12,
      top_p: 1,
      frequency_penalty: 0,
      presence_penalty: 0,
    };

    // request the openai api using axios
    await axios.post(openAIUrl, body, { headers }).then(async (response) => {
      // update the object with the competion result
      if (response.status == 200 && response.data.choices) {
        const completion = response.data.choices[0]!.text;
        obj.status!.completion = completion;
        console.log("completion object to be sent out", obj);
        await k8sApiMC.patchClusterCustomObjectStatus(
          COMPLETION_GROUP,
          COMPLETION_VERSION,
          COMPLETION_PLURAL,
          obj.metadata.name!,
          obj,
        );
      }
    });
  } catch (error) {
    console.log("error while reconciling object: ", error);
    // save the error to a json file
    const errorFile = `${obj.metadata.name}.json`;
    fs.writeFileSync(errorFile, JSON.stringify(error));
    console.log("wrote to file");
  }
}

async function main() {
  await watchResource();
}

function log(message: string) {
  console.log(`${new Date().toLocaleString()}: ${message}`);
}

main().catch((err) => {
  console.log(err);
  // write the error to a file
  fs.writeFile("error.json", err, function (e: any) {
    if (e) {
      console.log(e);
    }
    console.log("The file was saved!");
  });
});

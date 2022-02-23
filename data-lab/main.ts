import axios from "axios";
// import libraries
const path = require("path");
const fs = require("fs");

// list out the contents of the current directory
let localFiles = fs.readdirSync("./examples");

// return a map of {fileName => {path: path, content: content}}
const getFilesFromIssue = (issue) => {
  // extract a map of the files from the issue based on the following regex:
  // /`(@([a-zA-Z0-9_\-]+):(.+))`/g
  const fileRegex = /`(@([a-zA-Z0-9_\-]+):(.+))`/g;
  // create a map from the filename to the filepath
  const fileMap = new Map();
  // extract the string from group 3 of the regex
  let match;
  while ((match = fileRegex.exec(issue))) {
    let [name, path] = [match[2], match[3]];
    if (!fileMap.has(name)) {
      // define the file object here
      fileMap.set(name, {
        path: path,
        content: "",
        updatedContent: "",
      });
    } else {
      console.error(`duplicate file name ${name}`);
    }
  }
  return fileMap;
};

const getIssuesForDirectory = (dir) => {
  /* return a list of issue objects of the form: 
		{
			relevantFiles: string[],
			content: string,
			issueNumber: number
		}
	*/
  let issues = fs.readdirSync(dir).map((issueNo) => {
    let issuePath = path.join(dir, issueNo);
    let issue = fs.readFileSync(path.join(issuePath, "issue.md"), "utf8");
    return {
      relevantFiles: getFilesFromIssue(issue),
      content: issue,
      // convert the issueNo to a number
      issueNumber: parseInt(issueNo),
    };
  });
  // sort the issues by issue number
  issues.sort((a, b) => a.issueNumber - b.issueNumber);
  return issues;
};

const generateExamples = () => {
  examples = new Map();
  // loop through local files in ./examples
  // examples is a directory
  for (const dirname of fs.readdirSync("./examples")) {
    console.log(dirname);
    examples.set(dirname, {
      issues: [], // list of issues, sorted by their number
      files: new Map(), // map of {fileName => {path: path, content: content, updatedContent}}
    });

    let issueDir = path.join("./examples", dirname, "issues");
    examples.get(dirname).issues = getIssuesForDirectory(issueDir);
    if (examples.get(dirname).issues.length === 0) {
      // no issues, so skip this directory
      continue;
    }
    let firstIssue = examples.get(dirname).issues[0];
    console.log("first item in issues", firstIssue);
    // populate the files map
    examples.get(dirname).files = getFilesFromIssue(firstIssue.content);
    populateFiles(
      examples.get(dirname).files,
      path.join("./examples", dirname, "files")
    );
  }
  console.log(examples);
};

// @issue: string
// @files: map of {fileName: string => {path: string, content: string, updatedContent: string}}
const buildPrompt = (issue, files) => {
  // read the prefix from 'prefix.md'
  const prefix = fs.readFileSync("./prefix.md", "utf8");

  let prompt = `${prefix}
## Description of issues:
${issue}

## Original files:
`;

  i = 0;
  for (const [fileName, file] of files) {
    prompt += `# @${fileName}\n${file.content}\n`;
    // only place the delimiting string if in-between files
    if (files.size > 1 && i < files.size - 1) {
      prompt += "---\n";
    }
    i++;
  }

  prompt += `
## Updated files:
`;
  return prompt;
};

// process the completion & place it into the files' updatedContent field
const completionToFiles = (completion, files) => {
  let completions = completion.split("---");
  // go through the list of completions, extract the filename and set the updated content
  for (const cmpltn of completions) {
    // extract the file tag from the completion
    const fileTagRegex = /#\s*\@(.+)/g;
    const match = fileTagRegex.exec(cmpltn);
    if (match !== null) {
      const fileTag = match[1];
      if (files.has(fileTag)) {
        // find the line containing the fileTag and remove all lines up to and including the fileTag
        const lines = cmpltn.split("\n");
        let i = 0;
        for (const line of lines) {
          i++;
          if (line.includes(fileTag)) {
            break;
          }
        }
        // remove the lines from the completion
        const newCompletion = lines.slice(i).join("\n");
        // set the updated content
        files.get(fileTag).updatedContent = newCompletion;
      }
    }
  }
};

// to create the completion
const getCompletion = async (prompt, maxTokens, stopSequences) => {
  stopSequences = stopSequences || ["####"];
  const headers = {
    // get OPENAI_API_KEY from env
    Authorization: `Bearer ${process.env.OPENAI_API_KEY}`,
    "Content-Type": "application/json",
  };
  // console.log("headers", headers);
  const body = {
    prompt: prompt,
    max_tokens: maxTokens | 128,
    stop: stopSequences,
    temperature: 0.2,
    top_p: 1,
    frequency_penalty: 0,
    presence_penalty: 0,
  };
  let completion;

  // request the openai api using axios
  await axios.post(OPENAI_API_URL, body, { headers }).then(async (response) => {
    // update the object with the competion result
    if (response.status == 200 && response.data.choices) {
      if (response.data.choices.length > 0) {
        completion = response.data.choices[0].text;
      } else {
        console.error("no completion found");
      }
    }
  });
  return completion;
};

const updateFilesFromCompletion = async (files, prompt, maxTokens) => {
  const completion = await getCompletion(prompt, maxTokens, ["####"]);
  completionToFiles(completion, files);
};

// filesMap: map of {fileName: string => {path: string, content: string, updatedContent: string}}
// basePath: string
const writeFilesToCompletionsDir = (filesMap, basePath) => {
  // write the updated files into the completions directory using their same path as the original files
  for (const [_, file] of filesMap) {
    let outputPath = path.join(basePath, "completions", file.path);
    console.log("writing to: ", outputPath);
    fs.mkdirSync(path.dirname(outputPath), { recursive: true });
    // write the file and create parent directories, if needed
    fs.writeFileSync(outputPath, file.updatedContent);
  }
};

const generateCompletionsOnExample = async (basePath, example) => {
  for (let k = 0; k < example.issues.length; k++) {
    let issue = example.issues[k];
    let issuePath = path.join(
      basePath,
      "issues",
      example.issues[k].issueNumber.toString()
    );
    console.log(`::: example ${issue.issueNumber} in "${issuePath}" :::`);
    console.log(`\`\`\`\n${issue.content}\n\`\`\`\n`);

    // load relevant files from the given issue

    if (k > 0) {
      for (let [_, file] of example.files) {
        file.content = file.updatedContent;
        file.updatedContent = "";
      }
    }

    // build the initial prompt
    let initialPrompt = buildPrompt(example.issues[k].content, example.files);
    // write the prompt to a file
    fs.writeFileSync(path.join(issuePath, "prompt.md"), initialPrompt);

    // update the files with the completion
    let completion = await getCompletion(initialPrompt, 512, ["####"]);

    // write the completion to a file
    fs.writeFileSync(
      path.join(issuePath, "completion.md"),
      [initialPrompt, completion].join()
    );

    completionToFiles(completion, example.files);
    console.log("files after converting completions", example.files);
    writeFilesToCompletionsDir(example.files, issuePath);

    // update base files with the completion
    writeUpdatedContentToFiles(example.files, basePath);
  }
};

const main = async () => {
  // iterate through each example in the examples directory and generate completions
  let examples = fs.readdirSync("./examples");
  console.log(examples);
  for (const [dirname, example] of examples) {
    console.log(`===== processing example '${dirname}' ======`);
    // generateCompletionsOnExample(path.join('./examples', dirname), example);
  }
};

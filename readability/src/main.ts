import { JSDOM } from "jsdom";
import AWS, { Polly } from "aws-sdk";
import express from "express";
import { Readability, isProbablyReaderable } from "@mozilla/readability";

interface error {
  msg: string;
  detail?: string;
}

const bucket = process.env.RDR_BUCKET_NAME;
const prefix = process.env.RDR_PREFIX_NAME;
if (!bucket) {
  console.error("missing bucket name");
  process.exit(1);
}
if (!prefix) {
  console.error("missing prefix name");
  process.exit(1);
}

function parseDoc(docstring: Buffer | string, url: string) {
  const parser = new JSDOM(docstring, {
    url,
  });
  const reader = new Readability(parser.window.document);
  return reader.parse();
}

function checkReadability(data: Buffer, url: URL): boolean {
  const parsed = new JSDOM(data, { url: url.toString() });
  return isProbablyReaderable(parsed.window.document);
}

function fetchReadableArticle(data: Buffer, url: URL): string | undefined {
  const parsed = parseDoc(data, url.toString());
  if (!parsed) {
    return;
  }
  return parsed.textContent;
}

function getQuery(req: express.Request, prop: string): string | undefined {
  const val = req.query[prop];
  if (!val || typeof val !== "string") {
    return;
  }
  return val;
}

const app = express();
AWS.config.update({ region: "eu-west-2" });
const p = new Polly();

app.use((req, _, next) => {
  console.log(req.method, decodeURI(req.url));
  next();
});
app.use(express.raw({ type: "text/html", limit: "2mb" }));

app.get("/health", (_, res) => {
  res.status(200).send("ok");
});

app.post("/check", (req, res) => {
  const errors: error[] = [];
  const val = getQuery(req, "url");
  const url = new URL(val || "");
  if (!url) {
    errors.push({ msg: "missing url" });
    res.status(400).json({ errors });
    return;
  }
  const isReadable = checkReadability(req.body, new URL(url));
  console.log("is readable: ", isReadable);
  res
    .status(200)
    .setHeader("Content-Type", "application/json")
    .send({ is_readable: isReadable });
});

app.post("/extract", (req, res) => {
  const errors: error[] = [];
  const urlVal = getQuery(req, "url");
  if (!urlVal) {
    errors.push({ msg: "bad url parameter" });
    res.status(400).json({ errors });
    return;
  }
  const url = new URL(urlVal);
  if (!url) {
    res.status(400).json("missing url");
    return;
  }
  const parsed = fetchReadableArticle(req.body, new URL(url));
  if (!parsed) {
    errors.push({ msg: "failed to parse document" });
    res.status(400).json({ errors });
    return;
  }
  return res.status(200).setHeader("Content-Type", "text/html").send(parsed);
});

app.post("/synthesize", async (req, res) => {
  const errors: error[] = [];
  if (!Buffer.isBuffer(req.body)) {
    errors.push({ msg: "invalid body" });
    res.status(400).json({ errors });
    return;
  }
  const text = req.body.toString();

  const task = p.startSpeechSynthesisTask({
    OutputS3BucketName: bucket,
    OutputS3KeyPrefix: prefix,
    OutputFormat: "mp3",
    VoiceId: "Joanna",
    Engine: "standard",
    Text: text,
  });
  const output = (await task.promise()).SynthesisTask;
  if (!output) {
    errors.push({
      msg: "failed to synthesize",
      detail: "no output",
    });
    res.status(500).json({ errors });
    return;
  }
  switch (output.TaskStatus) {
    case "scheduled":
    case "inProgress":
    case "completed":
      res.status(202).json({
        jobId: output.TaskId,
        status: output.TaskStatus,
        reason: output.TaskStatusReason,
        errors,
      });
      return;
    default:
      errors.push({
        msg: "failed to synthesize",
        detail: output.TaskStatus,
      });
      res.status(500).json({ errors });
      return;
  }
});

app.get("/status", async (req, res) => {
  const errors: error[] = [];
  const jobId = getQuery(req, "jobId");
  if (!jobId) {
    errors.push({ msg: "missing jobId" });
    res.status(400).json({ errors });
    return;
  }

  const output = await p
    .getSpeechSynthesisTask({
      TaskId: jobId,
    })
    .promise();
  if (!output) {
    errors.push({ msg: "bad response when fetching job" });
    res.status(401).json({ errors });
    return;
  }

  const status = output.SynthesisTask?.TaskStatus;
  const reason = output.SynthesisTask?.TaskStatusReason;
  res.status(200).json({
    jobId,
    status,
    reason,
    errors,
  });
  return;
});

process.on("SIGINT", () => {
  process.exit(1);
});

app.listen(process.env.RDR_PORT || 5000, () => {
  console.log(`Server is running on ${process.env.RDR_PORT || "5000"}`);
});

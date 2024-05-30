import { JSDOM } from "jsdom";
import express from "express";
import { Readability, isProbablyReaderable } from "@mozilla/readability";

process.on("SIGINT", () => {
  process.exit(1);
});

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

function geturl(req: express.Request): URL | undefined {
  const url = req.query.url;
  if (!url || typeof url !== "string") {
    return;
  }
  return new URL(url);
}

const app = express();

app.use((req, _, next) => {
  console.log(req.method, decodeURI(req.url));
  next();
});
app.use(express.raw({ type: "text/html", limit: "2mb" }));

app.post("/check", (req, res) => {
  const url = geturl(req);
  if (!url) {
    res.status(400).send("missing url");
    return;
  }
  const isReadable = checkReadability(req.body, new URL(url));
  res
    .status(200)
    .setHeader("Content-Type", "application/json")
    .send({ is_readable: isReadable });
});

app.get("/health", (_, res) => {
  res.status(200).send("ok");
});

app.post("/check", (req, res) => {
  const url = geturl(req);
  if (!url) {
    res.status(400).send("missing url");
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
  const url = geturl(req);
  if (!url) {
    res.status(400).send("missing url");
    return;
  }
  const parsed = fetchReadableArticle(req.body, new URL(url));
  console.log("parsed", parsed);
  if (!parsed) {
    res.status(400).send("failed to parse");
    return;
  }
  return res.status(200).setHeader("Content-Type", "text/html").send(parsed);
});

app.listen(process.env.RDR_HOST, () => {
  console.log(`Server is running on ${process.env.RDR_HOST}`);
});

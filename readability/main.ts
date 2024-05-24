import { JSDOM } from "jsdom";
import express from "express";
import { Readability, isProbablyReaderable } from "@mozilla/readability";

function parseDoc(docstring: Buffer | string, url: string): any {
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
  if (!parsed?.textContent) {
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

app.get("/health", (_, res) => {
  res.status(200).send("ok");
});

app.use(express.raw({ type: "text/html" })).post("/check", (req, res) => {
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

app.post("/extract", (req, res) => {
  const url = geturl(req);
  if (!url) {
    res.status(400).send("missing url");
    return;
  }
  const parsed = fetchReadableArticle(req.body, new URL(url));
  if (!parsed) {
    res.status(400).send("failed to parse");
    return;
  }
  return res.status(200).setHeader("Content-Type", "text/html").send(parsed);
});

app.listen(5000, () => {
  console.log("Server is running on port 5000");
});

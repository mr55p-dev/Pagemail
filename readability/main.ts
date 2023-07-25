import { JSDOM } from "jsdom";
import { Readability } from "@mozilla/readability";

function parseDoc(docstring: Buffer, url: string): any {
  const parser = new JSDOM(docstring, {
    url,
  });
  const reader = new Readability(parser.window.document);
  return reader.parse();
}

function main() {
  const siteURI = process.argv[2];
  if (!siteURI) {
    throw new Error("Did not provide a site URI");
  }
  process.stdin.once("data", (data) => {
    const parsed = parseDoc(data, siteURI);
    process.stdout.write(JSON.stringify(parsed));
  });
}

main();

import { JSDOM } from "jsdom";
import { Readability } from "@mozilla/readability";

function parseDoc(docstring: Buffer | string, url: string): any {
  const parser = new JSDOM(docstring, {
    url,
  });
  const reader = new Readability(parser.window.document);
  return reader.parse();
}

async function main() {
  const siteURI = process.argv[2];
  const url = new URL(siteURI)
  if (!url) {
    console.error("Did not provide a site URI");
	process.exit(1)
  }

  const res = await fetch(siteURI)
  if (!res.ok) {
	console.error("Failed to fetch")
	process.exit(2)
  }

  const body = await res.text()
  const parsed = parseDoc(body, url.toString());
  if (!parsed?.textContent) {
	process.exit(3)
  }

  process.stdout.write(JSON.stringify(parsed));
  process.exit(0)
}

main()

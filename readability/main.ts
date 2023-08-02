import { parseArgs } from "node:util";
import { JSDOM } from "jsdom";
import { Readability, isProbablyReaderable } from "@mozilla/readability";

function parseDoc(docstring: Buffer | string, url: string): any {
  const parser = new JSDOM(docstring, {
    url,
  });
  const reader = new Readability(parser.window.document);
  return reader.parse();
}

async function fetchQuick(data: Buffer, url: URL) {
  const parsed = new JSDOM(data, { url: url.toString() });
  const isReadable = isProbablyReaderable(parsed.window.document);
  if (isReadable) {
    process.exit(0);
  } else {
    process.exit(1);
  }
}

async function fetchReadableArticle(data: Buffer, url: URL) {
  const parsed = parseDoc(data, url.toString());
  if (!parsed?.textContent) {
    process.exit(3);
  }

  process.stdout.write(JSON.stringify(parsed));
  process.exit(0);
}

async function main() {
  const { values } = parseArgs({
    options: {
      check: {
        type: "boolean",
        default: false,
      },
      url: {
        type: "string",
      },
    },
    strict: true,
  });
  const url = new URL(values.url || "");
  if (!url) {
    console.error("Did not provide a site URI");
    process.exit(1);
  }

  process.stdin.addListener("data", (stream) => {
    if (values.check) {
      fetchQuick(stream, url);
    } else {
      fetchReadableArticle(stream, url);
    }
  });
}

main();

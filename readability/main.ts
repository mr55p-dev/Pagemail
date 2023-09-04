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
    console.log("checks passed");
    process.exit(0);
  } else {
    console.error("checks failed");
    process.exit(1);
  }
}

async function fetchReadableArticle(data: Buffer, url: URL) {
  const parsed = parseDoc(data, url.toString());
  if (!parsed?.textContent) {
    process.exit(3);
  }

  const out = JSON.stringify(parsed)
  process.stdout.write(out);
}

class Message {
  ready = false;
  clHeaderSz = 4;
  cl = 0;
  buf: Buffer | undefined;
  bufIdx = 0;
  done = false;
  onExit: VoidFunction | undefined;
  frIdx = 0;

  initialize(headerBytes: Buffer): number {
    if (this.ready) {
      return 0;
    }

    // Set the content length and buffer
    this.cl = headerBytes.readUIntBE(0, this.clHeaderSz);
    this.buf = Buffer.allocUnsafe(this.cl);
    this.ready = true;
    return this.clHeaderSz;
  }

  processFrame(data: Buffer, offset?: number) {
    let bufOffset = offset || 0;

    if (!this.buf) {
      throw new Error("initialized and no buffer present");
    }

    for (; bufOffset < data.length; bufOffset++) {
      const char = data.at(bufOffset);
      if (char === undefined) {
        throw new Error(`Undefined character of data at index ${bufOffset}`);
      }

      this.buf.writeUInt8(char, this.bufIdx);
      this.bufIdx++;
      if (this.bufIdx === this.cl) {
        this.done = true;
		return true
      }
    }
    this.frIdx++;
  }

  attachStream(callback: (buf: Buffer) => void) {
    process.stdin.addListener("readable", () => {
      const data = process.stdin.read();
      if (!data || this.done) {
        process.stdin.destroy();
        return;
      }

      let offset;
      if (!this.ready) {
        offset = this.initialize(data.subarray(0, this.clHeaderSz));
      }

      this.processFrame(data, offset);
      process.stdin.emit("end");
	  if (this.done) {
		callback(this.buf!)
	  }
    });
  }
}

function main() {
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

  if (!values.url) {
    console.error("Did not provide a site URI");
    process.exit(1);
  }

  try {
    var url = new URL(values.url || "");
  } catch (_) {
    console.error("URL is invalid");
    process.exit(1);
  }



  const msg = new Message()
  msg.attachStream(buf => {
	if (values.check) {
	  fetchQuick(buf, url)
	} else {
	  fetchReadableArticle(buf, url)
	}
  })
}

main();

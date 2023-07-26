"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const jsdom_1 = require("jsdom");
const readability_1 = require("@mozilla/readability");
function parseDoc(docstring, url) {
    const parser = new jsdom_1.JSDOM(docstring, {
        url,
    });
    const reader = new readability_1.Readability(parser.window.document);
    return reader.parse();
}
function main() {
    const siteURI = process.argv[2];
    if (!siteURI) {
        throw new Error("Did not provide a site URI");
    }
    process.stdin.once("data", (data) => {
        const parsed = parseDoc(data, siteURI);
        if (!(parsed === null || parsed === void 0 ? void 0 : parsed.textContent)) {
            process.exit(1);
        }
        process.stdout.write(JSON.stringify(parsed));
    });
}
main();

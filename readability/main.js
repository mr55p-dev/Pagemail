"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var jsdom_1 = require("jsdom");
var express_1 = require("express");
var readability_1 = require("@mozilla/readability");
function parseDoc(docstring, url) {
    var parser = new jsdom_1.JSDOM(docstring, {
        url: url,
    });
    var reader = new readability_1.Readability(parser.window.document);
    return reader.parse();
}
function checkReadability(data, url) {
    var parsed = new jsdom_1.JSDOM(data, { url: url.toString() });
    return (0, readability_1.isProbablyReaderable)(parsed.window.document);
}
function fetchReadableArticle(data, url) {
    var parsed = parseDoc(data, url.toString());
    if (!(parsed === null || parsed === void 0 ? void 0 : parsed.textContent)) {
        return;
    }
    return parsed.textContent;
}
function geturl(req) {
    var url = req.query.url;
    if (!url || typeof url !== "string") {
        return;
    }
    return new URL(url);
}
var app = (0, express_1.default)();
app.get("/health", function (_, res) {
    res.status(200).send("ok");
});
app.use(express_1.default.raw({ type: "text/html" })).post("/check", function (req, res) {
    var url = geturl(req);
    if (!url) {
        res.status(400).send("missing url");
        return;
    }
    var isReadable = checkReadability(req.body, new URL(url));
    res
        .status(200)
        .setHeader("Content-Type", "application/json")
        .send({ is_readable: isReadable });
});
app.post("/extract", function (req, res) {
    var url = geturl(req);
    if (!url) {
        res.status(400).send("missing url");
        return;
    }
    var parsed = fetchReadableArticle(req.body, new URL(url));
    if (!parsed) {
        res.status(400).send("failed to parse");
        return;
    }
    return res.status(200).setHeader("Content-Type", "text/html").send(parsed);
});
app.listen(5000, function () {
    console.log("Server is running on port 5000");
});

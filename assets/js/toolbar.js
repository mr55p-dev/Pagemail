function pasteContents(nodeId) {
    navigator.clipboard.readText().then((contents) => {
        document.getElementById(nodeId).value = contents;
    });
}

function copyText(text) {
    navigator.clipboard
        .writeText(text)
        .then(() => console.log("Copied"))
        .catch((err) => console.error("Failed to copy", err));
}

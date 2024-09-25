function pasteContents(nodeId) {
    navigator.clipboard.readText().then((contents) => {
        document.getElementById(nodeId).value = contents;
    });
}

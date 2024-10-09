function pasteContents(nodeId) {
    navigator.clipboard.readText().then((contents) => {
        document.getElementById(nodeId).value = contents;
    });
}

function copyText(button) {
    const parent = button.closest("article");
    const [span] = button.getElementsByTagName("span");
    navigator.clipboard
        .writeText(parent.dataset.url)
        .then(() => {
            button.dataset.copied = true;
            const originalText = span.textContent;
            span.textContent = "Copied";
            setTimeout(() => {
                button.dataset.copied = false;
                span.textContent = originalText;
            }, 1500);
        })
        .catch((err) => {
            button.dataset.failed = true
            console.error("Failed to copy", err)
        });
}

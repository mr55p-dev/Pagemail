export default function ({ setFieldCallback }) {
    const pasteFromClipboard = () => {
        navigator.clipboard.readText().then((text) => {
            setFieldCallback(text)
        })
    }

    // Check if the navigator api supports clipboard

    return navigator.clipboard.readText && (
        <div>
            <button onClick={pasteFromClipboard}>Paste from clipboard</button>
        </div>
    )
}
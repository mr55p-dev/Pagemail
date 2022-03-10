import { useEffect, useState } from "react"

export default function PasteFromClipboard ({ setFieldCallback }) {
    const pasteFromClipboard = () => navigator.clipboard.readText().then(text => setFieldCallback(text))

    return (navigator.clipboard) && (
        <button className="border-2 submit-enabled py-2
        md:col-span-2 md:row-span-1 md:row-start-2 md:col-start-11 my-1" onClick={pasteFromClipboard}>Paste from clipboard</button>
    )
}
import { usePageMetadata } from "../lib/hooks"
import { ICard } from "../lib/typeAliases"

export default function PageCard({ data, deleteCallback }) {

    const metadata = true && data?.metadata
    const title = metadata?.title || data.url
    const desc = metadata?.description || ""
    const date = new Date(1000 * data.timeAdded.seconds).toLocaleDateString();
    const url = new URL(data.url);

    return (
        <div className="border-2 border-sky-700 bg-sky-50">
            <h3 className="">{title}</h3>
            <p className="">{desc}</p>
            <a className="" href={url.toString()} target="_blank" rel="noreferrer">{url.hostname}</a>
            <small className="">{date}</small>
            <button className="" onClick={() => {deleteCallback(data.id)}}>X</button>
        </div>
    )
}
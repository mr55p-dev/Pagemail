import Link from "next/link"
import { usePageMetadata } from "../lib/hooks"
import { ICard } from "../lib/typeAliases"

export default function PageCard({ data, deleteCallback }) {

    const metadata = true && data?.metadata
    const title = metadata?.title || data.url
    const desc = metadata?.description || "No description available"
    const date = new Date(1000 * data.timeAdded.seconds).toLocaleDateString();
    const url = new URL(data.url);

    return (
        <div className="border-2 rounded shadow-sm border-sky-700 bg-sky-50 p-2 flex flex-col justify-between">
            <div className="">
                <h3 className="text-lg font-semibold">{title}</h3>
                <p className="">{desc}</p>
            </div>
            <div className="">
                <div className="text-center grid grid-rows-2 gap-1 sm:grid-cols-2 sm:grid-rows-1 sm:space-x-2">
                    <a className="underline border-2 border-sky-700 btn-colour rounded p-3 text-sky-700 md:grow" href={url.toString()} target="_blank" rel="noreferrer">Open {url.hostname}</a>
                    <button className="hover:bg-red-700 hover:text-sky-50 border-2 border-red-700 transition-colors rounded p-3 md:grow" onClick={() => {deleteCallback(data.id)}}>Remove</button>
                </div>
                <small className="font-light">Saved {date}</small>
            </div>
        </div>
    )
}
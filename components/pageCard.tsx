import Link from "next/link"
import { usePageMetadata } from "../lib/hooks"
import { ICard } from "../lib/typeAliases"

export default function PageCard({ data, deleteCallback }) {

    const metadata = true && data?.metadata
    const title = metadata?.title || data.url
    const desc = metadata?.description || "No description available"
    const date = new Date(1000 * data.timeAdded.seconds).toLocaleDateString();
    const url = new URL(data.url);

    let sitename = url.hostname
    if (sitename.startsWith("www.")) {
        sitename = sitename.replace("www.", "")
    }

    return (
        <div className="text-secondary dark:text-secondary-dark bg-primary dark:bg-primary-dark
                        grid grid-rows-6 grid-cols-1 sm:grid-rows-1 sm:grid-cols-12 items-stretch
                        odd:border-y-2 last:border-b-2 border-secondary dark:border-secondary-dark p-2 md:py-4">
            <a href={url.toString()} target="_blank" rel="noreferrer"
                            className="row-span-4 sm:row-span-2 sm:col-span-10 sm:p-2 lg:col-span-11
                            whitespace-nowrap overflow-hidden sm:mr-2 hover:cursor-pointer group">
                <h3 className="transition-all font-semibold group-hover:underline underline-offset-2 text-xl
                    group-hover:text-tertiary group-hover:font-bold">{title}</h3>
                <p className="mb-2">{(desc.length > 150) ? desc.substr(0, 150-3) + '...' : desc}</p>
                <small className="font-light">Saved {date}</small>
            </a>
            <div className="row-span-2 sm:col-span-2 lg:col-span-1 flex items-center justify-around
                        border-tertiary border-2 text-semibold text-md
                        hover:bg-tertiary hover:text-primary hover:dark:text-primary-dark
                        text-secondary dark:text-secondary-dark transition-all">
                <button className="w-full h-full" onClick={() => {deleteCallback(data.id)}}>Remove</button>
            </div>
        </div>
    )
}
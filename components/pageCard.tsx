import { usePageMetadata } from "../lib/hooks"

export default function PageCard(props) {
    const metadata = usePageMetadata(props.url, props.token)

    const title = metadata?.title || props.url
    const desc = metadata?.description || ""
    const url = new URL(props.url);
    const date = props.dateCreated;

    return (
        <div className="pages-item">
            <h3 className="pages-card-title">{title}</h3>
            <p className="pages-card-description">{desc}</p>
            <a className="pages-card-url" href={url.toString()}>{url.hostname}</a>
            <small className="pages-card-date">{date}</small>
            <button className="btn pages-card-btn" onClick={() => {props.deleteCallback(props.documentID)}}>X</button>
        </div>
    )
}
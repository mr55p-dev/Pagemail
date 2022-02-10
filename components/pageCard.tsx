export default function PageCard(props) {
    return(
        <div className="pages-item">
            <h3 className="pages-card-title">{props.url}</h3>
            <a className="pages-card-url" href={props.url}>Open</a>
            <small className="pages-card-date">{props.dateCreated}</small>
            <button className="btn pages-card-btn" onClick={() => {props.deleteCallback(props.documentID)}}>X</button>
        </div>
    )
}
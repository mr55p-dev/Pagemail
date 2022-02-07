export default function PageCard(props) {
    return(
        <div className="pages-item">
            <h3 className="page-card-title">{props.url}</h3>
            <a className="page-card-url" href={props.url}>Open</a>
            <small className="page-card-date">{props.dateCreated}</small>
            <button className="btn pages-card-btn" onClick={() => {props.deleteCallback(props.documentID)}}>X</button>
        </div>
    )
}
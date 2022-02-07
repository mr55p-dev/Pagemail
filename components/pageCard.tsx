export default function PageCard(props) {
    return(
        <div className="pages-item">
            <h3 className="page-card-title">{props.title}</h3>
            <small><a className="page-card-url" href={props.url}>{props.url}</a></small>
            <small>{props.documentID}</small>
            <button onClick={() => {props.deleteCallback(props.documentID)}}>X</button>
        </div>
    )
}
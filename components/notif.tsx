import ReactDOM from "react-dom";
import { useRendered } from "../lib/hooks";
import { INotifProp } from "../lib/typeAliases";

export default function Notif({ show, state }: INotifProp): JSX.Element {
    const rendered = useRendered()
    const content = rendered ? (
        <div className="notif-root">
            <div className={state?.style}>
                <em className="notif-title">{state?.title}</em>
                <small className="notif-text">{state?.text}</small>
            </div>
        </div>
        ) : null;

    return (show) ? ReactDOM.createPortal(
        content,
        document.getElementById("notif-root")
    ) : null;
}
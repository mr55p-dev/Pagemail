import ReactDOM from "react-dom";
import { useEffect, useState } from "react";

export default function Modal({ children, show, onClose }) {
    const [isBrowser, setIsBrowser] = useState(false);

    useEffect(() => {
        setIsBrowser(true);
    }, []);

    const handleClose = (e) => {
        e.preventDefault();
        onClose();
    }


    const modalContent = show ? (
        <div className="modal-container">
            {children}
            <button onClick={handleClose}>X</button>
        </div>
    ) : null;

    if (isBrowser) {
        return ReactDOM.createPortal(
            modalContent,
            document.getElementById("modal-root")
        );
      } else {
        return null;
      }
}
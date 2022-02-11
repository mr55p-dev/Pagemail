import ReactDOM from "react-dom";
import { useEffect, useState } from "react";
import { useRendered } from "../lib/hooks";

export default function Modal({ children, show, onClose }) {
    const rendered = useRendered();

    const handleClose: React.UIEventHandler<HTMLButtonElement> = (e) => {
        e.preventDefault();
        onClose();
    }


    const modalContent = show ? (
        <div className="modal-container">
            {children}
            <button onClick={handleClose}>X</button>
        </div>
    ) : null;

    if (rendered) {
        return ReactDOM.createPortal(
            modalContent,
            document.getElementById("modal-root")
        );
      } else {
        return null;
      }
}
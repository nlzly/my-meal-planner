import { ReactNode, useState } from "react";
import { createPortal } from "react-dom";
import "./Modal.css"

type ModalProps ={
    children : ReactNode,
    text : string
}

export default function Modal({children, text}: ModalProps) {
    const [showModal, setShowModal] = useState(false);
    const handleOverlayClick = (e: React.MouseEvent<HTMLDivElement>) => {
        if (e.target === e.currentTarget){
            setShowModal(false);
        }
    }
    return (
        <>
            <button onClick={() => setShowModal(true)}>
                {text}
            </button>
            {showModal && createPortal(
                <div className="modal-overlay" onClick={handleOverlayClick}>
                    <div className="modal-content">
                        {children}
                    </div>
                </div>,
                document.body
            )}
        </>
    )
}
import { ReactNode, useState } from "react";
import { createPortal } from "react-dom";
import "./Modal.css"

type ModalProps ={
    children : ReactNode
    isOpen: boolean;
    onClose: () => void;
}

export default function Modal({children, isOpen, onClose}: ModalProps) {
    if(!isOpen) return null;

    const handleOverlayClick = (e: React.MouseEvent<HTMLDivElement>) => {
        if (e.target === e.currentTarget){
            onClose();
        }
    }
    return createPortal(
        <div className="modal-overlay" onClick={handleOverlayClick}>
          <div className="modal-content">
            {children}
            <button onClick={onClose}>Close</button>
          </div>
        </div>,
        document.body
      );
}
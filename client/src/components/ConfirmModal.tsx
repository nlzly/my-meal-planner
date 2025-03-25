import React from 'react';
import Modal from './Modal';

interface ConfirmModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  title: string;
  message: string;
}

const ConfirmModal: React.FC<ConfirmModalProps> = ({
  isOpen,
  onClose,
  onConfirm,
  title,
  message
}) => {
  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <div className="confirm-modal">
        <h2>{title}</h2>
        <p>{message}</p>
        <div className="confirm-modal-buttons">
          <button className="cancel-button" onClick={onClose}>
            Cancel
          </button>
          <button className="confirm-button" onClick={onConfirm}>
            Confirm
          </button>
        </div>
      </div>
    </Modal>
  );
};

export default ConfirmModal; 
import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const CancelMoveConfirmationModal = ({ onClose, onSubmit, moveID, title, content, submitText, closeText }) => (
  <Modal>
    <ModalClose handleClick={() => onClose()} />
    <ModalTitle>
      <h3>{title}</h3>
    </ModalTitle>
    <p>{content}</p>
    <ModalActions autofocus="true">
      <Button
        data-focus="true"
        className="usa-button--destructive"
        type="submit"
        data-testid="modalSubmitButton"
        onClick={() => onSubmit(moveID)}
      >
        {submitText}
      </Button>
      <Button className="usa-button--secondary" type="button" onClick={() => onClose()} data-testid="modalBackButton">
        {closeText}
      </Button>
    </ModalActions>
  </Modal>
);

CancelMoveConfirmationModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,

  title: PropTypes.string,
  content: PropTypes.string,
  submitText: PropTypes.string,
  closeText: PropTypes.string,
};

CancelMoveConfirmationModal.defaultProps = {
  title: 'Are you sure?',
  content:
    'If you proceed with cancellation and later want to process a move request for these orders, you will have to start again with a new move.',
  submitText: 'Cancel move',
  closeText: 'Keep move',
};

CancelMoveConfirmationModal.displayName = 'CancelMoveConfirmationModal';

export default connectModal(CancelMoveConfirmationModal);

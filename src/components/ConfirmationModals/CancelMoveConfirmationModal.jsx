import React from 'react';
import { Button } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const CancelMoveConfirmationModal = ({ onClose, onSubmit, moveId, title, content, submitText, closeText }) => (
  <Modal>
    <ModalClose handleClick={() => onClose()} />
    <ModalTitle>
      <h3 data-testid="modaltitle">{title}</h3>
    </ModalTitle>
    <p>{content}</p>
    <ModalActions autofocus="true">
      <Button data-focus="true" className="usa-button--destructive" type="submit" onClick={() => onSubmit(moveId)}>
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

  moveId: PropTypes.string.isRequired,

  title: PropTypes.string,
  content: PropTypes.string,
  submitText: PropTypes.string,
  closeText: PropTypes.string,
};

CancelMoveConfirmationModal.defaultProps = {
  title: 'Are you sure?',
  content:
    'You’ll lose all the information in this move. If you want it back later, you’ll have to request a new move.',
  submitText: 'Cancel move',
  closeText: 'Keep move',
};

CancelMoveConfirmationModal.displayName = 'CancelMoveConfirmationModal';

export default connectModal(CancelMoveConfirmationModal);

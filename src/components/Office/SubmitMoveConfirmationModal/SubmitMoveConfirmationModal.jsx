import React from 'react';
import PropTypes from 'prop-types';
import { Button, Overlay, ModalContainer } from '@trussworks/react-uswds';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const SubmitMoveConfirmationModal = ({ onClose, onSubmit }) => (
  <div>
    <Overlay />
    <ModalContainer>
      <Modal>
        <ModalClose handleClick={() => onClose()} />
        <ModalTitle>
          <h3>Are you sure?</h3>
        </ModalTitle>
        <p>You can’t make changes after you submit the move.</p>
        <ModalActions>
          <Button className="usa-button--submit" type="submit" onClick={() => onSubmit()}>
            Yes, submit
          </Button>
          <Button
            className="usa-button--tertiary"
            type="button"
            onClick={() => onClose()}
            data-testid="modalBackButton"
          >
            Cancel
          </Button>
        </ModalActions>
      </Modal>
    </ModalContainer>
  </div>
);

SubmitMoveConfirmationModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

SubmitMoveConfirmationModal.displayName = 'SubmitMoveConfirmationModal';

export default connectModal(SubmitMoveConfirmationModal);

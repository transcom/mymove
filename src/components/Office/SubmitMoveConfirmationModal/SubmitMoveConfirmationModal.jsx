import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { connectModal, ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';

export const SubmitMoveConfirmationModal = ({ onClose, onSubmit, bodyText }) => (
  <div data-testid="SubmitMoveConfirmationModal">
    <Overlay />
    <ModalContainer>
      <Modal>
        <ModalClose handleClick={() => onClose()} />
        <ModalTitle>
          <h2>Are you sure?</h2>
        </ModalTitle>
        <p>{bodyText}</p>
        <ModalActions>
          <Button className="usa-button--submit" type="submit" onClick={() => onSubmit()}>
            Yes, submit
          </Button>
          <Button
            className="usa-button--tertiary"
            type="button"
            onClick={() => onClose()}
            data-testid="modalCancelButton"
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
  bodyText: PropTypes.string,
};

SubmitMoveConfirmationModal.defaultProps = {
  bodyText: "You can't make changes after you submit the move.",
};

SubmitMoveConfirmationModal.displayName = 'SubmitMoveConfirmationModal';

export default connectModal(SubmitMoveConfirmationModal);

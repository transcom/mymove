import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { connectModal, ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';

export const SubmitMoveConfirmationModal = ({ onClose, onSubmit, isShipment }) => (
  <div data-testid="SubmitMoveConfirmationModal">
    <Overlay />
    <ModalContainer>
      <Modal onClose={() => onClose()}>
        <ModalClose handleClick={() => onClose()} />
        <ModalTitle>
          <h2>Are you sure?</h2>
        </ModalTitle>
        <p>You can’t make changes after you submit the {isShipment ? 'shipment' : 'move'}.</p>
        <ModalActions>
          <Button
            className="usa-button--secondary"
            type="button"
            onClick={() => onClose()}
            data-testid="modalCancelButton"
          >
            Cancel
          </Button>
          <Button className="usa-button--submit" type="submit" onClick={() => onSubmit()}>
            Yes, submit
          </Button>
        </ModalActions>
      </Modal>
    </ModalContainer>
  </div>
);

SubmitMoveConfirmationModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  isShipment: PropTypes.bool,
};

SubmitMoveConfirmationModal.defaultProps = {
  isShipment: false,
};

SubmitMoveConfirmationModal.displayName = 'SubmitMoveConfirmationModal';

export default connectModal(SubmitMoveConfirmationModal);

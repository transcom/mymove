import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { connectModal, ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';

export const CompletePPMCloseoutConfirmationModal = ({ onClose, onSubmit }) => (
  <div data-testid="CompletePPMCloseoutConfirmationModal">
    <Overlay />
    <ModalContainer>
      <Modal>
        <ModalClose handleClick={() => onClose()} data-testid="modalCloseButton" />
        <ModalTitle>
          <h2>Are you sure you want to complete the PPM Review?</h2>
        </ModalTitle>
        <ModalActions>
          <Button
            className="usa-button--tertiary"
            type="button"
            onClick={() => onClose()}
            data-testid="modalBackButton"
          >
            No
          </Button>
          <Button className="usa-button--submit" type="submit" onClick={() => onSubmit()}>
            Yes
          </Button>
        </ModalActions>
      </Modal>
    </ModalContainer>
  </div>
);

CompletePPMCloseoutConfirmationModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

CompletePPMCloseoutConfirmationModal.displayName = 'CompletePPMCloseoutConfirmationModal';

export default connectModal(CompletePPMCloseoutConfirmationModal);

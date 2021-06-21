import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';
import { Overlay, ModalContainer } from 'components/MigratedModal/MigratedModal';

export const SCRequestShipmentCancellationModal = ({ onClose, onSubmit, shipmentID }) => (
  <div>
    <Overlay />
    <ModalContainer>
      <Modal>
        <ModalClose handleClick={() => onClose()} />
        <ModalTitle>
          <h3>Are you sure?</h3>
        </ModalTitle>
        <p>
          You’ll lose all the information in this shipment. If you want it back later, you’ll have to request a new
          shipment.
        </p>
        <ModalActions>
          <Button className="usa-button--destructive" type="submit" onClick={() => onSubmit(shipmentID)}>
            Delete shipment
          </Button>
          <Button
            className="usa-button--tertiary"
            type="button"
            onClick={() => onClose()}
            data-testid="modalBackButton"
          >
            Keep shipment
          </Button>
        </ModalActions>
      </Modal>
    </ModalContainer>
  </div>
);

SCRequestShipmentCancellationModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,

  shipmentID: PropTypes.string.isRequired,
};

SCRequestShipmentCancellationModal.displayName = 'SCRequestShipmentCancellationModal';

export default connectModal(SCRequestShipmentCancellationModal);

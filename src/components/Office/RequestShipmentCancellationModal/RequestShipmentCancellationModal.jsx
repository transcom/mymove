import React from 'react';
import PropTypes from 'prop-types';
import { Button, Overlay, ModalContainer } from '@trussworks/react-uswds';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const RequestShipmentCancellationModal = ({ onClose, onSubmit, shipmentInfo }) => (
  <div>
    <Overlay />
    <ModalContainer>
      <Modal>
        <ModalClose handleClick={() => onClose()} />
        <ModalTitle>
          <h3>Request shipment cancellation</h3>
        </ModalTitle>
        <p>
          Movers will be notified that this shipment should be canceled. They will confirm or deny this request based on
          whether or not service items have been charged to the shipment yet.
        </p>
        <ModalActions>
          <Button
            className="usa-button--tertiary"
            type="button"
            onClick={() => onClose()}
            data-testid="modalBackButton"
          >
            Back
          </Button>
          <Button
            className="usa-button--destructive"
            type="submit"
            onClick={() => onSubmit(shipmentInfo.moveTaskOrderID, shipmentInfo.shipmentID, shipmentInfo.ifMatchEtag)}
          >
            Request Cancellation
          </Button>
        </ModalActions>
      </Modal>
    </ModalContainer>
  </div>
);

RequestShipmentCancellationModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  shipmentInfo: PropTypes.shape({
    shipmentID: PropTypes.string.isRequired,
    ifMatchEtag: PropTypes.string.isRequired,
    moveTaskOrderID: PropTypes.string.isRequired,
  }).isRequired,
};

RequestShipmentCancellationModal.displayName = 'RequestShipmentCancellationModal';

export default connectModal(RequestShipmentCancellationModal);

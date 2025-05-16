import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';

const RequestShipmentCancellationModal = ({ onClose, onSubmit, shipmentInfo }) => (
  <div>
    <Overlay />
    <ModalContainer>
      <Modal onClose={() => onClose()}>
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
            className="usa-button--secondary"
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

export default RequestShipmentCancellationModal;

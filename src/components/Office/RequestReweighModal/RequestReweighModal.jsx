import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';

const RequestReweighModal = ({ onClose, onSubmit, shipmentInfo }) => (
  <div>
    <Overlay />
    <ModalContainer onClose={() => onClose()}>
      <Modal>
        <ModalClose handleClick={() => onClose()} />
        <ModalTitle>
          <h3>Request a reweigh</h3>
        </ModalTitle>
        <p>
          This will notify the movers of the request. They&apos;ll reweigh the shipment if it&apos;s still possible.
        </p>
        <ModalActions>
          <Button
            className="usa-button--secondary"
            type="button"
            onClick={() => onClose()}
            data-testid="modalBackButton"
          >
            Cancel
          </Button>
          <Button
            className="usa-button"
            type="submit"
            onClick={() => onSubmit(shipmentInfo.shipmentID, shipmentInfo.ifMatchEtag)}
          >
            Submit request
          </Button>
        </ModalActions>
      </Modal>
    </ModalContainer>
  </div>
);

RequestReweighModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  shipmentInfo: PropTypes.shape({
    shipmentID: PropTypes.string.isRequired,
    ifMatchEtag: PropTypes.string.isRequired,
  }).isRequired,
};

RequestReweighModal.displayName = 'RequestReweighModal';

export default RequestReweighModal;

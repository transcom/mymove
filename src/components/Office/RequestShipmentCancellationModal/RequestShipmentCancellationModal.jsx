import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './RequestShipmentCancellationModal.module.scss';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const RequestShipmentCancellationModal = ({ onClose, onSubmit }) => (
  <Modal>
    <ModalClose className={styles.closeButton} handleClick={onClose} />
    <ModalTitle>
      <h3>Request shipment cancellation</h3>
    </ModalTitle>
    <p>
      Movers will be notified that this shipment should be canceled. They will confirm or deny this request based on
      whether or not service items have been charged to the shipment yet.
    </p>
    <ModalActions>
      <Button secondary className={styles.cancelButton} type="button" onClick={onClose} data-testid="modalBackButton">
        Back
      </Button>
      <Button className={styles.requestButton} type="submit" onClick={onSubmit}>
        Request Cancellation
      </Button>
    </ModalActions>
  </Modal>
);

RequestShipmentCancellationModal.propTypes = {
  onClose: PropTypes.func,
  onSubmit: PropTypes.func,
};

RequestShipmentCancellationModal.defaultProps = {
  onClose: () => {},
  onSubmit: () => {},
};

RequestShipmentCancellationModal.displayName = 'RequestShipmentCancellationModal';

export default connectModal(RequestShipmentCancellationModal);

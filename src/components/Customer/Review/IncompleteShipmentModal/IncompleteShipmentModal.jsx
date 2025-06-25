import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from 'components/Customer/Review/IncompleteShipmentModal/IncompleteShipmentModal.module.scss';
import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const IncompleteShipmentModal = ({ closeModal, shipmentLabel, shipmentMoveCode, shipmentType }) => (
  <Modal className={styles.Modal} onClose={closeModal}>
    <ModalClose handleClick={closeModal} />
    <ModalTitle>
      <h3>INCOMPLETE SHIPMENT</h3>
    </ModalTitle>
    <p>
      <b>
        {shipmentLabel}: #{shipmentMoveCode}
      </b>
    </p>
    <p>
      You have not finished adding all the details required for your {shipmentType} shipment. Click <b>Edit</b> to
      review your {shipmentType} information, add any missing information, then proceed to submit the request.
    </p>
    <ModalActions>
      <Button secondary type="button" onClick={closeModal} className={styles.Button}>
        OK
      </Button>
    </ModalActions>
  </Modal>
);

IncompleteShipmentModal.propTypes = {
  closeModal: PropTypes.func.isRequired,
  shipmentLabel: PropTypes.string.isRequired,
  shipmentMoveCode: PropTypes.string.isRequired,
  shipmentType: PropTypes.string.isRequired,
};

IncompleteShipmentModal.displayName = 'IncompleteShipmentModal';

export default connectModal(IncompleteShipmentModal);

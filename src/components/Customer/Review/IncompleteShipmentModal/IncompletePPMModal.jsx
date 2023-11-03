import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from 'components/Customer/Review/IncompleteShipmentModal/IncompletePPMModal.module.scss';
import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const IncompletePPMModal = ({ closeModal, data }) => (
  <Modal className={styles.Modal}>
    <ModalClose handleClick={closeModal} />
    <ModalTitle>
      <h3>INCOMPLETE SHIPMENT</h3>
    </ModalTitle>
    <p>
      {JSON.parse(data).shipmentLabel}: {JSON.parse(data).shipmentIDAbbrevLabel}
    </p>
    <p>You have elected not to use advanced request....blah blah</p>
    <ModalActions>
      <Button secondary type="button" onClick={closeModal} className={styles.Button}>
        OK
      </Button>
    </ModalActions>
  </Modal>
);

IncompletePPMModal.propTypes = {
  closeModal: PropTypes.func,
  data: PropTypes.string,
};

IncompletePPMModal.defaultProps = {
  closeModal: () => {},
  data: null,
};

IncompletePPMModal.displayName = 'IncompletePPMModal';

export default connectModal(IncompletePPMModal);

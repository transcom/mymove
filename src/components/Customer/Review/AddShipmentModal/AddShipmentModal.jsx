import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from 'components/Customer/Review/AddShipmentModal/AddShipmentModal.module.scss';
import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const AddShipmentModal = ({ closeModal, enablePPM, enableNTS, enableNTSR }) => (
  <Modal className={styles.Modal} onClose={closeModal}>
    <ModalClose handleClick={closeModal} />
    <ModalTitle>
      <h3>Reasons you might need another shipment</h3>
    </ModalTitle>
    {enablePPM && enableNTS && enableNTSR ? (
      <>
        <ul>
          <li>
            You plan to have an <strong>HHG</strong> and a <strong>PPM (DITY)</strong> — you want the government to pay
            professional movers, and you also want to be reimbursed for moving some things yourself.
          </li>
          <li>You have additional belongings to move from or to a very different location, like another city.</li>
          <li>You need to schedule another type of shipment, like an NTS. This would be listed on your orders.</li>
        </ul>
        <p>If none of these apply to you, you probably don&apos;t need another shipment.</p>
      </>
    ) : (
      <p>
        Some shipment types are still being developed and will become available at a later date. If you require a
        different shipment type not yet available in MilMove, contact your transportation office for further
        instructions.
      </p>
    )}
    <ModalActions>
      <Button secondary type="button" onClick={closeModal} className={styles.Button}>
        Got it
      </Button>
    </ModalActions>
  </Modal>
);

AddShipmentModal.propTypes = {
  closeModal: PropTypes.func,
  enablePPM: PropTypes.bool,
  enableNTS: PropTypes.bool,
  enableNTSR: PropTypes.bool,
};

AddShipmentModal.defaultProps = {
  closeModal: () => {},
  enablePPM: true,
  enableNTS: true,
  enableNTSR: true,
};

AddShipmentModal.displayName = 'AddShipmentModal';

export default connectModal(AddShipmentModal);

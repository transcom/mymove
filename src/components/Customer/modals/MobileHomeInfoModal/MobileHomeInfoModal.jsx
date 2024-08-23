import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const MobileHomeInfoModal = ({ closeModal }) => (
  <Modal>
    <ModalClose handleClick={closeModal} />
    <ModalTitle>
      <h3>Boat & Mobile homes info</h3>
    </ModalTitle>
    <h4>
      <strong>Mobile Home shipment</strong>
    </h4>
    <p>This option is for privately owned mobile homes.</p>
    <ModalActions>
      <Button secondary type="button" onClick={closeModal}>
        Got it
      </Button>
    </ModalActions>
  </Modal>
);

MobileHomeInfoModal.propTypes = {
  closeModal: PropTypes.func,
};

MobileHomeInfoModal.defaultProps = {
  closeModal: () => {},
};

MobileHomeInfoModal.displayName = 'MobileHomeInfoModal';

export default connectModal(MobileHomeInfoModal);

import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './ErrorModal.module.scss';
import Modal, { ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';
import SystemError from 'components/SystemError';

export const ErrorModal = ({ closeModal, errorMessage, displayHelpDeskLink }) => (
  <Modal className={styles.Modal}>
    <ModalClose handleClick={closeModal} />
    <SystemError>
      {errorMessage}
      {displayHelpDeskLink && <a href="mailto:usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil">Technical Help Desk</a>}
    </SystemError>
    <ModalActions>
      <Button secondary type="button" onClick={closeModal} className={styles.Button}>
        ok
      </Button>
    </ModalActions>
  </Modal>
);

ErrorModal.propTypes = {
  closeModal: PropTypes.func.isRequired,
};

ErrorModal.displayName = 'ErrorModal';

export default connectModal(ErrorModal);

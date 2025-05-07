import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './IncorrectXlsxErrorModal.module.scss';

import Modal, { ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';
import SystemError from 'components/SystemError';

export const IncorrectXlsxErrorModal = ({ closeModal, errorMessage }) => (
  <Modal className={styles.Modal}>
    <ModalClose handleClick={closeModal} />
    <SystemError>{errorMessage}</SystemError>
    <ModalActions>
      <Button secondary type="button" onClick={closeModal} className={styles.Button}>
        ok
      </Button>
    </ModalActions>
  </Modal>
);

IncorrectXlsxErrorModal.propTypes = {
  closeModal: PropTypes.func.isRequired,
};

IncorrectXlsxErrorModal.displayName = 'IncorrectXlsxErrorModal';

export default connectModal(IncorrectXlsxErrorModal);

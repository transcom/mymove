import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './DownloadAOAErrorModal.module.scss';
import Modal, { ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';
import SystemError from 'components/SystemError';

export const DownloadAOAErrorModal = ({ closeModal }) => (
  <Modal className={styles.Modal}>
    <ModalClose handleClick={closeModal} />
    <SystemError>
      Something went wrong downloading PPM AOA paperwork. Please try again later. If that doesn&apos;t fix it, contact
      the &nbsp;<a href="mailto:usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@mail.mil">Technical Help Desk</a>.
    </SystemError>
    <ModalActions>
      <Button secondary type="button" onClick={closeModal} className={styles.Button}>
        ok
      </Button>
    </ModalActions>
  </Modal>
);

DownloadAOAErrorModal.propTypes = {
  closeModal: PropTypes.func.isRequired,
};

DownloadAOAErrorModal.displayName = 'DownloadAOAErrorModal';

export default connectModal(DownloadAOAErrorModal);

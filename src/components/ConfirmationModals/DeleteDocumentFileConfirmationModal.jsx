import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import bytes from 'bytes';
import moment from 'moment';

import styles from './DeleteDocumentFileConfirmationModal.module.scss';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const DeleteDocumentFileConfirmationModal = ({ closeModal, submitModal, fileInfo }) => (
  <Modal onClose={closeModal}>
    <ModalClose handleClick={closeModal} />
    <ModalTitle>
      <h3>Are you sure you want to delete this file?</h3>
    </ModalTitle>
    <div className={styles.fileInfo}>
      <p className={styles.fileName}>{fileInfo.filename}</p>
      <p className={styles.fileSizeAndTime}>
        <span className={styles.uploadFileSize}>{bytes(fileInfo.bytes)}</span>
        <span>Uploaded {moment(fileInfo.createdAt).format('DD MMM YYYY h:mm A')}</span>
      </p>
    </div>
    <ModalActions>
      <Button className="usa-button--secondary" type="button" onClick={closeModal} data-testid="cancel-delete">
        No, keep it
      </Button>
      <Button
        data-testid="confirm-delete"
        data-focus="true"
        className="usa-button--destructive"
        type="submit"
        onClick={submitModal}
      >
        Yes, delete
      </Button>
    </ModalActions>
  </Modal>
);

DeleteDocumentFileConfirmationModal.propTypes = {
  closeModal: PropTypes.func.isRequired,
  submitModal: PropTypes.func.isRequired,
  fileInfo: PropTypes.object.isRequired,
};

DeleteDocumentFileConfirmationModal.displayName = 'DeleteDocumentFileConfirmationModal';

export default connectModal(DeleteDocumentFileConfirmationModal);

import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const DeleteEvaluationReportConfirmationModal = ({ closeModal, submitModal, isDeleteFromTable }) => (
  <Modal onClose={closeModal}>
    <ModalClose handleClick={closeModal} />
    <ModalTitle>
      <h3>Are you sure you want to {isDeleteFromTable ? 'delete' : 'cancel'} this report?</h3>
    </ModalTitle>
    <p>You cannot undo this action.</p>
    <ModalActions>
      <Button className="usa-button--secondary" type="button" onClick={closeModal} data-testid="modalBackButton">
        No, keep it
      </Button>
      <Button autoFocus data-focus="true" className="usa-button--destructive" type="submit" onClick={submitModal}>
        Yes, {isDeleteFromTable ? 'delete' : 'cancel'}
      </Button>
    </ModalActions>
  </Modal>
);

DeleteEvaluationReportConfirmationModal.propTypes = {
  closeModal: PropTypes.func.isRequired,
  submitModal: PropTypes.func.isRequired,
  isDeleteFromTable: PropTypes.bool,
};

DeleteEvaluationReportConfirmationModal.defaultProps = {
  isDeleteFromTable: false,
};

DeleteEvaluationReportConfirmationModal.displayName = 'DeleteEvaluationReportConfirmationModal';

export default connectModal(DeleteEvaluationReportConfirmationModal);

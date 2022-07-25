import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const DeleteEvaluationReportConfirmationModal = ({ closeModal, submitModal }) => (
  <Modal>
    <ModalClose handleClick={closeModal} />
    <ModalTitle>
      <h3>Are you sure you want to cancel this report?</h3>
    </ModalTitle>
    <p>You cannot undo this action.</p>
    <ModalActions autofocus="true">
      <Button data-focus="true" className="usa-button--destructive" type="submit" onClick={submitModal}>
        Yes, Cancel
      </Button>
      <Button className="usa-button--secondary" type="button" onClick={closeModal} data-testid="modalBackButton">
        No, keep it
      </Button>
    </ModalActions>
  </Modal>
);

DeleteEvaluationReportConfirmationModal.propTypes = {
  closeModal: PropTypes.func.isRequired,
  submitModal: PropTypes.func.isRequired,
};

DeleteEvaluationReportConfirmationModal.displayName = 'DeleteEvaluationReportConfirmationModal';

export default connectModal(DeleteEvaluationReportConfirmationModal);

import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const BulkAssignmentModal = ({ onClose, onSubmit, title, content, submitText, closeText }) => (
  <Modal>
    <ModalClose handleClick={() => onClose()} />
    <ModalTitle>
      <h3>{title}</h3>
    </ModalTitle>
    <p>{content}</p>
    <ModalActions autofocus="true">
      <Button
        data-focus="true"
        className="usa-button--destructive"
        type="submit"
        data-testid="modalSubmitButton"
        onClick={() => onSubmit()}
        // onSubmit({
        //   queueType: 'COUNSELING',
        //   saveBulkAssignmentData: {
        //     userData: [
        //       {
        //         userId: '1cc7ce3a-c4ef-4582-a58a-e70665b412b0',
        //         moveAssignments: 2,
        //       },
        //     ],
        //     moveData: ['b3baf6ce-f43b-437c-85be-e1145c0ddb96', 'fee7f916-35a6-4c0b-9ea6-a1d8094b3ed3'],
        //   },
        // })
        // }
      >
        {submitText}
      </Button>
      <Button className="usa-button--secondary" type="button" onClick={() => onClose()} data-testid="modalBackButton">
        {closeText}
      </Button>
    </ModalActions>
  </Modal>
);

BulkAssignmentModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,

  title: PropTypes.string,
  content: PropTypes.string,
  submitText: PropTypes.string,
  closeText: PropTypes.string,
};

BulkAssignmentModal.defaultProps = {
  title: 'Bulk Assignment',
  content: 'Here we will display moves to be assigned in bulk.',
  submitText: 'Save',
  closeText: 'Cancel',
};

BulkAssignmentModal.displayName = 'BulkAssignmentModal';

export default connectModal(BulkAssignmentModal);

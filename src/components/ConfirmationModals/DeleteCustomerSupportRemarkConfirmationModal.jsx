import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const DeleteCustomerSupportRemarkConfirmationModal = ({
  onClose,
  onSubmit,
  customerSupportRemarkID,
  title,
  content,
  submitText,
  closeText,
}) => (
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
        onClick={() => onSubmit(customerSupportRemarkID)}
      >
        {submitText}
      </Button>
      <Button className="usa-button--secondary" type="button" onClick={() => onClose()} data-testid="modalBackButton">
        {closeText}
      </Button>
    </ModalActions>
  </Modal>
);

DeleteCustomerSupportRemarkConfirmationModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,

  customerSupportRemarkID: PropTypes.string.isRequired,

  title: PropTypes.string,
  content: PropTypes.string,
  submitText: PropTypes.string,
  closeText: PropTypes.string,
};

DeleteCustomerSupportRemarkConfirmationModal.defaultProps = {
  title: 'Are you sure you want to delete this remark?',
  content: 'You cannot undo this action.',
  submitText: 'Yes, Delete',
  closeText: 'No, keep it',
};

DeleteCustomerSupportRemarkConfirmationModal.displayName = 'DeleteCustomerSupportRemarkConfirmationModal';

export default connectModal(DeleteCustomerSupportRemarkConfirmationModal);

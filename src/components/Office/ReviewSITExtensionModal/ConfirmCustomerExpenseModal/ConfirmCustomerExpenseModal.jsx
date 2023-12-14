import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from '../ReviewSITExtensionModal.module.scss';

import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';

const ConfirmCustomerExpenseModal = ({ onClose, onSubmit }) => {
  return (
    <div>
      <Modal>
        <ModalClose handleClick={() => onClose()} />
        <ModalTitle>
          <h2>Convert to Customer Expense</h2>
        </ModalTitle>
        <div className={styles.ModalPanel}>
          <p>Are you sure that you would like to convert to Customer Expense?</p>
          <ModalActions>
            <Button data-testid="convertToCustomerExpenseConfirmationYes" onClick={() => onSubmit()}>
              Yes
            </Button>
            <Button data-testid="convertToCustomerExpenseConfirmationNo" onClick={() => onClose()}>
              No
            </Button>
          </ModalActions>
        </div>
      </Modal>
    </div>
  );
};

ConfirmCustomerExpenseModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};
export default ConfirmCustomerExpenseModal;

import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from '../ReviewSITExtensionModal.module.scss';

import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';

const ConfirmCustomerExpenseModal = ({ setShowConfirmModal, values, setValues }) => {
  const handleConfirmYes = () => {
    setValues({
      ...values,
      convertToCustomerExpense: true,
    });
    setShowConfirmModal(false);
  };
  const handleConfirmNo = () => {
    setValues({
      ...values,
      convertToCustomerExpense: false,
    });
    setShowConfirmModal(false);
  };

  return (
    <Modal onClose={handleConfirmNo}>
      <ModalClose handleClick={handleConfirmNo} />
      <ModalTitle>
        <h2>Convert to Customer Expense</h2>
      </ModalTitle>
      <div className={styles.ModalPanelNoBorder}>
        <p>Are you sure that you would like to convert to Customer Expense?</p>
        <ModalActions>
          <Button
            type="button"
            secondary
            data-testid="convertToCustomerExpenseConfirmationNo"
            onClick={handleConfirmNo}
          >
            No
          </Button>
          <Button type="button" data-testid="convertToCustomerExpenseConfirmationYes" onClick={handleConfirmYes}>
            Yes
          </Button>
        </ModalActions>
      </div>
    </Modal>
  );
};

ConfirmCustomerExpenseModal.propTypes = {
  setShowConfirmModal: PropTypes.func.isRequired,
  values: PropTypes.object.isRequired,
  setValues: PropTypes.func.isRequired,
};
export default ConfirmCustomerExpenseModal;

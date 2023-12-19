import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from '../ReviewSITExtensionModal.module.scss';

import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';

const ConfirmCustomersExpenseModal = ({ onClose, onSubmit }) => {
  return (
    <Modal>
      <ModalClose handleClick={onClose} />
      <ModalTitle>
        <h2>Convert to Customer&apos;s Expense</h2>
      </ModalTitle>
      <div className={styles.ModalPanel}>
        <p>Are you sure that you would like to convert to Customer&apos;s Expense?</p>
        <ModalActions>
          <Button data-testid="convertToCustomersExpenseConfirmationYes" onClick={onSubmit}>
            Yes
          </Button>
          <Button data-testid="convertToCustomersExpenseConfirmationNo" onClick={onClose}>
            No
          </Button>
        </ModalActions>
      </div>
    </Modal>
  );
};

ConfirmCustomersExpenseModal.propTypes = {
  onSubmit: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
};
export default ConfirmCustomersExpenseModal;

import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from '../ReviewSITExtensionModal.module.scss';

import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';

const ConfirmMembersExpenseModal = ({ onClose, onSubmit }) => {
  return (
    <Modal>
      <ModalClose handleClick={onClose} />
      <ModalTitle>
        <h2>Convert to Member&apos;s Expense</h2>
      </ModalTitle>
      <div className={styles.ModalPanel}>
        <p>Are you sure that you would like to convert to Member&apos;s Expense?</p>
        <ModalActions>
          <Button data-testid="convertToMembersExpenseConfirmationYes" onClick={onSubmit}>
            Yes
          </Button>
          <Button data-testid="convertToMembersExpenseConfirmationNo" onClick={onClose}>
            No
          </Button>
        </ModalActions>
      </div>
    </Modal>
  );
};

ConfirmMembersExpenseModal.propTypes = {
  onSubmit: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
};
export default ConfirmMembersExpenseModal;

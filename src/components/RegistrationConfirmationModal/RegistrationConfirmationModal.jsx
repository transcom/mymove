import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './RegistrationConfirmationModal.module.scss';

import Modal, { ModalActions, ModalTitle, connectModal } from 'components/Modal/Modal';

export const RegistrationConfirmationModal = ({ onSubmit }) => {
  return (
    <Modal>
      <ModalTitle className={styles.center}>
        <h3>Registration Confirmation</h3>
      </ModalTitle>
      <p className={styles.center}>
        Your MilMove & Okta accounts have successfully been created. <br />
        It is required that you first sign-in with a Common Access Card (CAC). <br />
        <br />
        You will now be redirected to the Okta sign-in page where you will click on the &quot;Sign in with PIV / CAC
        Card&quot; button to sign in.
        <br />
      </p>
      <ModalActions autofocus="true">
        <Button data-focus="true" type="submit" data-testid="modalSubmitButton" onClick={() => onSubmit()}>
          Continue
        </Button>
      </ModalActions>
    </Modal>
  );
};

RegistrationConfirmationModal.propTypes = {
  onSubmit: PropTypes.func.isRequired,
};

RegistrationConfirmationModal.displayName = 'RegistrationConfirmationModal';

export default connectModal(RegistrationConfirmationModal);

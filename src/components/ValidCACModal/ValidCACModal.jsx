import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from '../../pages/CreateAccount/CreateAccount.module.scss';

import smartCard from 'shared/images/smart-card.png';
import Modal, { ModalTitle, ModalActions, connectModal } from 'components/Modal/Modal';

export const ValidCACModal = ({ onClose, onSubmit }) => (
  <Modal onClose={onClose}>
    <ModalTitle className={styles.center}>
      <h3>Do you have a valid CAC?</h3>
    </ModalTitle>
    <p className={styles.center}>
      <img src={smartCard} width="200" height="200" alt="" />
    </p>
    <p className={styles.center}>
      Common Access Card (CAC) authentication is required at first sign-in. <br />
      If you do not have a CAC, do not request your account here. <br />
      You must visit your nearest personal property office where they will assist you with creating your MilMove
      account.
    </p>
    <ModalActions autofocus="true">
      <Button className="usa-button--secondary" type="button" onClick={() => onClose()} data-testid="modalBackButton">
        No
      </Button>
      <Button data-focus="true" type="submit" data-testid="modalSubmitButton" onClick={() => onSubmit()}>
        Yes
      </Button>
    </ModalActions>
  </Modal>
);

ValidCACModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

ValidCACModal.displayName = 'ValidCACModal';

export default connectModal(ValidCACModal);

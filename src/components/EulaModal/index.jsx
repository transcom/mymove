import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './EulaModal.module.scss';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const EulaModal = ({ closeModal, acceptTerms }) => (
  <Modal className={styles.eulaModal}>
    <ModalClose handleClick={closeModal} />
    <ModalTitle>
      <h3>Standard mandatory DoD Notice and consent Banner</h3>
    </ModalTitle>
    <p>
      You are accessing a U.S. Government (USG) Information System (IS) that is provided for USG-authorized use only.
    </p>
    <p>By using this IS (which includes any device attached to this IS), you consent to the following conditions:</p>
    <ul>
      <li>
        The USG routinely intercepts and monitors communications on this IS for purposes including, but not limited to,
        penetration testing, COMSEC monitoring, network operations and defense, personnel misconduct (PM), law
        enforcement (LE), and counterintelligence (CI) investigations.
      </li>
      <li>At any time, the USG may inspect and seize data stored on this IS.</li>
      <li>
        Communications using, or data stored on, this IS are not private, are subject to routine monitoring,
        interception, and search, and may be disclosed or used for any USG-authorized purpose.
      </li>
      <li>
        This IS includes security measures (e.g., authentication and access controls) to protect USG interests--not for
        your personal benefit or privacy.
      </li>
      <li>
        Notwithstanding the above, using this IS does not constitute consent to PM, LE or CI investigative searching or
        monitoring of the content of privileged communications, or work product, related to personal representation or
        services by attorneys, psychotherapists, or clergy, and their assistants. Such communications and work product
        are private and confidential. See User Agreement for details.
      </li>
    </ul>
    <ModalActions>
      <Button aria-label="Accept EULA Terms" type="button" onClick={acceptTerms}>
        Accept terms
      </Button>
      <Button aria-label="Cancel" secondary type="button" onClick={closeModal}>
        Cancel
      </Button>
    </ModalActions>
  </Modal>
);

EulaModal.propTypes = {
  acceptTerms: PropTypes.func.isRequired,
  closeModal: PropTypes.func.isRequired,
};

EulaModal.displayName = 'EulaModal';

export default connectModal(EulaModal);

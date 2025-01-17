import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './BulkAssignmentModal.module.scss';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const BulkAssignmentModal = ({ onClose, onSubmit, title, submitText, closeText, bulkAssignmentData }) => (
  <div>
    <Modal className={styles.BulkModal}>
      <ModalClose handleClick={() => onClose()} />
      <ModalTitle>
        <h3>
          {title} (
          {bulkAssignmentData.bulkAssignmentMoveIDs == null ? 0 : bulkAssignmentData.bulkAssignmentMoveIDs.length})
        </h3>
      </ModalTitle>
      <div className={styles.BulkAssignmentTable}>
        <table>
          <tr>
            <th>User</th>
            <th>Workload</th>
            <th>Assignment</th>
          </tr>
          {bulkAssignmentData?.availableOfficeUsers?.map((user) => {
            return (
              <tr key={user}>
                <td>
                  <p data-testid="bulkAssignmentUser">
                    {user.lastName}, {user.firstName}
                  </p>
                </td>
                <td className={styles.BulkAssignmentDataCenter}>
                  <p data-testid="bulkAssignmentUserWorkload">{user.workload || 0}</p>
                </td>
                <td className={styles.BulkAssignmentDataCenter}>
                  <input className={styles.BulkAssignmentAssignment} type="number" />
                </td>
              </tr>
            );
          })}
        </table>
      </div>
      <ModalActions autofocus="true">
        <Button
          data-focus="true"
          className="usa-button--destructive"
          type="submit"
          data-testid="modalSubmitButton"
          onClick={() => onSubmit()}
        >
          {submitText}
        </Button>
        <Button className="usa-button--secondary" type="button" onClick={() => onClose()} data-testid="modalBackButton">
          {closeText}
        </Button>
      </ModalActions>
    </Modal>
  </div>
);

BulkAssignmentModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,

  title: PropTypes.string,
  submitText: PropTypes.string,
  closeText: PropTypes.string,
};

BulkAssignmentModal.defaultProps = {
  title: 'Bulk Assignment',
  submitText: 'Save',
  closeText: 'Cancel',
};

BulkAssignmentModal.displayName = 'BulkAssignmentModal';

export default connectModal(BulkAssignmentModal);

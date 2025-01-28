import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './BulkAssignmentModal.module.scss';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';
import { getBulkAssignmentData } from 'services/ghcApi';
import { milmoveLogger } from 'utils/milmoveLog';

export const BulkAssignmentModal = ({ onClose, onSubmit, title, submitText, closeText, queueType }) => {
  const [bulkAssignmentData, setBulkAssignmentData] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        getBulkAssignmentData(queueType).then((data) => {
          setBulkAssignmentData(data);
        });
      } catch (err) {
        milmoveLogger.error('Error fetching bulk assignment data:', err);
      }
    };
    fetchData();
  }, [queueType]);
  return (
    <div>
      <Modal className={styles.BulkModal}>
        <ModalClose handleClick={() => onClose()} />
        <ModalTitle>
          <h3>
            {title} ({bulkAssignmentData == null ? 0 : bulkAssignmentData.bulkAssignmentMoveIDs.length})
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
            className="usa-button--submit"
            type="submit"
            data-testid="modalSubmitButton"
            onClick={() => onSubmit()}
          >
            {submitText}
          </Button>
          <button className={styles.backButton} type="button" onClick={() => onClose()} data-testid="modalBackButton">
            {closeText}
          </button>
        </ModalActions>
      </Modal>
    </div>
  );
};
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

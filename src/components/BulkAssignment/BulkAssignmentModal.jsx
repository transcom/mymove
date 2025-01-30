import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './BulkAssignmentModal.module.scss';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';
import { getBulkAssignmentData } from 'services/ghcApi';
import { milmoveLogger } from 'utils/milmoveLog';
import { userName } from 'utils/formatters';

export const BulkAssignmentModal = ({ onClose, onSubmit, title, submitText, closeText, queueType }) => {
  const [bulkAssignmentData, setBulkAssignmentData] = useState(null);
  const [isDisabled, setIsDisabled] = useState(false);
  const [numberOfMoves, setNumberOfMoves] = useState(0);
  useEffect(() => {
    const fetchData = async () => {
      try {
        getBulkAssignmentData(queueType).then((data) => {
          setBulkAssignmentData(data);
          if (data.bulkAssignmentMoveIDs === undefined) {
            setIsDisabled(true);
            setNumberOfMoves(0);
          } else {
            setNumberOfMoves(data.bulkAssignmentMoveIDs.length);
          }
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
        <ModalClose handleClick={onClose} />
        <ModalTitle>
          <h3>
            {title} ({numberOfMoves})
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
                    <p data-testid="bulkAssignmentUser">{userName(user)}</p>
                  </td>
                  <td className={styles.BulkAssignmentDataCenter}>
                    <p data-testid="bulkAssignmentUserWorkload">{user.workload || 0}</p>
                  </td>
                  <td className={styles.BulkAssignmentDataCenter}>
                    <input className={styles.BulkAssignmentAssignment} type="number" min="0" />
                  </td>
                </tr>
              );
            })}
          </table>
        </div>
        <ModalActions autofocus="true">
          <Button
            disabled={isDisabled}
            data-focus="true"
            type="submit"
            data-testid="modalSubmitButton"
            onClick={() => onSubmit()}
          >
            {submitText}
          </Button>
          <Button secondary type="button" onClick={onClose} data-testid="modalCancelButton">
            {closeText}
          </Button>
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

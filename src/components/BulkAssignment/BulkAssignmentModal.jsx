import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './BulkAssignmentModal.module.scss';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

const data = {
  availableOfficeUsers: [
    {
      firstName: 'John',
      lastName: 'Snow',
      officeUserId: '123',
      workload: 0,
    },
    {
      firstName: 'Jane',
      lastName: 'Doe',
      officeUserId: '456',
      workload: 1,
    },
    {
      firstName: 'Jimmy',
      lastName: 'Page',
      officeUserId: '789',
      workload: 2,
    },
    {
      firstName: 'Peter',
      lastName: 'Pan',
      officeUserId: '101',
      workload: 3,
    },
    {
      firstName: 'Ringo',
      lastName: 'Starr',
      officeUserId: '111',
      workload: 4,
    },
    {
      firstName: 'George',
      lastName: 'Harrison',
      officeUserId: '121',
      workload: 5,
    },
    {
      firstName: 'Stuart',
      lastName: 'Skinner',
      officeUserId: '131',
      workload: 6,
    },
  ],
  bulkAssignmentMoveIDs: ['1', '2', '3', '4', '5'],
};

export const BulkAssignmentModal = ({
  onClose,
  onSubmit,
  title,
  submitText,
  closeText,
  //  bulkAssignmentData,
}) => (
  <div>
    <Modal className={styles.BulkModal}>
      <ModalClose handleClick={() => onClose()} />
      <ModalTitle>
        <h3>{title}</h3>
      </ModalTitle>
      <div className={styles.BulkAssignmentTable}>
        <table>
          <tr>
            <th>Select/Deselect All </th>
            <th>User</th>
            <th>Workload</th>
            <th>Assignment</th>
          </tr>
          {data.availableOfficeUsers.map((user) => {
            return (
              <tr>
                <td className={styles.BulkAssignmentDataCenter}>
                  <input type="checkbox" />
                </td>
                <td>
                  {user.firstName},{user.lastName}
                </td>
                <td className={styles.BulkAssignmentDataCenter}>{user.workload}</td>
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

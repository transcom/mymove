import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';

import styles from './BulkAssignmentModal.module.scss';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';
import { Form } from 'components/form';

const initialValues = {
  userData: [
    // {
    //   userId: '045c3048-df9a-4d44-88ed-8cd6e2100e08',
    //   moveAssignments: 2,
    // },
    // {
    //   userId: '4b1f2722-b0bf-4b16-b8c4-49b4e49ba42a',
    //   moveAssignments: 1,
    // },
  ],
  moveData: [],
};

export const BulkAssignmentModal = ({ onClose, onSubmit, title, submitText, closeText, bulkAssignmentData }) => (
  <Modal>
    <ModalClose handleClick={() => onClose()} />
    <ModalTitle>
      <h3>
        {title} (
        {bulkAssignmentData.bulkAssignmentMoveIDs == null ? 0 : bulkAssignmentData.bulkAssignmentMoveIDs.length})
      </h3>
    </ModalTitle>
    <div className={styles.BulkAssignmentTable}>
      <Formik
        onSubmit={(e) => {
          console.log(e);
          return onSubmit(e);
        }}
        initialValues={initialValues}
      >
        <Form>
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
                    <input
                      className={styles.BulkAssignmentAssignment}
                      type="number"
                      id="officeUserAssignment"
                      key={user}
                    />
                  </td>
                </tr>
              );
            })}
          </table>
          <ModalActions autofocus="true">
            <Button
              data-focus="true"
              className="usa-button--destructive"
              type="submit"
              data-testid="modalSubmitButton"
              // onClick={() =>
              //   onSubmit('COUNSELING', {
              //     userData: [
              //       {
              //         userId: '045c3048-df9a-4d44-88ed-8cd6e2100e08',
              //         moveAssignments: 2,
              //       },
              //       {
              //         userId: '4b1f2722-b0bf-4b16-b8c4-49b4e49ba42a',
              //         moveAssignments: 1,
              //       },
              //     ],
              //     moveData: [
              //       'b3baf6ce-f43b-437c-85be-e1145c0ddb96',
              //       'fee7f916-35a6-4c0b-9ea6-a1d8094b3ed3',
              //       '962ce8d2-03a2-435c-94ca-6b9ef6c226c1',
              //     ],
              //   })
              // }
            >
              {submitText}
            </Button>
            <Button
              className="usa-button--secondary"
              type="button"
              onClick={() => onClose()}
              data-testid="modalBackButton"
            >
              {closeText}
            </Button>
          </ModalActions>
        </Form>
      </Formik>
    </div>
  </Modal>
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

import React, { useEffect } from 'react';
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

export const BulkAssignmentModal = ({ onClose, onSubmit, title, submitText, closeText, bulkAssignmentData }) => {
  // adds move data to the initialValues obj
  useEffect(() => {
    console.log(initialValues);
    initialValues.moveData = bulkAssignmentData.bulkAssignmentMoveIDs;
    console.log(initialValues);
  });

  return (
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
          onSubmit={(values) => {
            onSubmit({
              bulkAssignmentSavePayload: values,
            });
          }}
          initialValues={initialValues}
        >
          {({ setValues, values }) => {
            const addUserAssignment = (user) => {
              let newUserAssignment;
              if (user.target.value !== '') {
                newUserAssignment = {
                  userId: user.target.id,
                  moveAssignments: user.target.value,
                };
              } else {
                newUserAssignment = {
                  userId: user.target.id,
                  moveAssignments: 0,
                };
              }

              setValues({
                ...values,
                userData: [newUserAssignment],
              });
            };

            return (
              <Form>
                <table>
                  <tr>
                    <th>User</th>
                    <th>Workload</th>
                    <th>Assignment</th>
                  </tr>
                  {bulkAssignmentData?.availableOfficeUsers?.map((user) => {
                    return (
                      <tr key={user.officeUserId}>
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
                            id={user.officeUserId}
                            onChange={addUserAssignment}
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
                    onClick={(bulkAssignmentSavePayload) => onSubmit(bulkAssignmentSavePayload)}
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
            );
          }}
        </Formik>
      </div>
    </Modal>
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

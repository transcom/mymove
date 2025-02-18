import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import * as Yup from 'yup';

import styles from './BulkAssignmentModal.module.scss';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';
import { Form } from 'components/form';
import { getBulkAssignmentData } from 'services/ghcApi';
import { milmoveLogger } from 'utils/milmoveLog';

const initialValues = {
  userData: [],
  moveData: [],
};

export const BulkAssignmentModal = ({ onClose, onSubmit, title, submitText, closeText, queueType }) => {
  const [isError, setIsError] = useState(false);
  const [bulkAssignmentData, setBulkAssignmentData] = useState(null);
  const [isDisabled, setIsDisabled] = useState(false);
  const [numberOfMoves, setNumberOfMoves] = useState(0);

  const errorMessage = 'Cannot assign more moves than are available.';

  const initUserData = (availableOfficeUsers) => {
    const officeUsers = [];
    availableOfficeUsers.forEach((user) => {
      const newUserAssignment = {
        ID: user.officeUserId,
        moveAssignments: 0,
      };
      officeUsers.push(newUserAssignment);
    });
    initialValues.userData = officeUsers;
  };

  useEffect(() => {
    const fetchData = async () => {
      try {
        getBulkAssignmentData(queueType).then((data) => {
          setBulkAssignmentData(data);
          initUserData(data?.availableOfficeUsers);

          if (data.bulkAssignmentMoveIDs === undefined) {
            setIsDisabled(true);
            setNumberOfMoves(0);
          } else {
            setNumberOfMoves(data.bulkAssignmentMoveIDs.length);
          }
        });
      } catch (err) {
        setBulkAssignmentData({});
        milmoveLogger.error('Error fetching bulk assignment data:', err);
      }
    };

    fetchData();
  }, [queueType]);

  // adds move data to the initialValues obj
  initialValues.moveData = bulkAssignmentData?.bulkAssignmentMoveIDs;

  const validationSchema = Yup.object().shape({
    assignment: Yup.number().min(0).typeError('Assignment must be a number'),
  });

  return (
    <div>
      <Modal>
        <ModalClose handleClick={() => onClose()} />
        <ModalTitle>
          <h3>
            {title} ({numberOfMoves})
          </h3>
        </ModalTitle>
        <div className={styles.BulkAssignmentTable}>
          <Formik
            onSubmit={(values) => {
              const totalAssignment = values?.userData?.reduce((sum, item) => sum + item.moveAssignments, 0);

              if (totalAssignment > numberOfMoves) {
                setIsError(true);
                return;
              }

              const bulkAssignmentSavePayload = values;
              onSubmit({ bulkAssignmentSavePayload });
              onClose();
            }}
            validationSchema={validationSchema}
            initialValues={initialValues}
          >
            {({ handleChange, setValues, values }) => {
              const handleAssignmentChange = (event, i) => {
                handleChange(event);
                setIsError(false);

                let newUserAssignment;
                if (event.target.value !== '') {
                  newUserAssignment = {
                    ID: event.target.id,
                    moveAssignments: +event.target.value,
                  };
                } else {
                  newUserAssignment = {
                    ID: event.target.id,
                    moveAssignments: 0,
                  };
                }

                const newValues = values;
                newValues.userData[i] = newUserAssignment;

                setValues({
                  ...values,
                  userData: newValues.userData,
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
                    {bulkAssignmentData?.availableOfficeUsers?.map((user, i) => {
                      return (
                        <tr key={user.officeUserId}>
                          <td>
                            <p data-testid="bulkAssignmentUser" className={styles.officeUserFormattedName}>
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
                              name="assignment"
                              id={user.officeUserId}
                              data-testid="assignment"
                              defaultValue={0}
                              min={0}
                              onChange={(event) => handleAssignmentChange(event, i)}
                            />
                          </td>
                        </tr>
                      );
                    })}
                  </table>
                  <ModalActions autofocus="true">
                    <Button disabled={isDisabled} data-focus="true" type="submit" data-testid="modalSubmitButton">
                      {submitText}
                    </Button>
                    <Button
                      type="button"
                      className={styles.button}
                      unstyled
                      onClick={() => onClose()}
                      data-testid="modalCancelButton"
                    >
                      {closeText}
                    </Button>
                  </ModalActions>
                  {isError && <div className={styles.errorMessage}>{errorMessage}</div>}
                </Form>
              );
            }}
          </Formik>
        </div>
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

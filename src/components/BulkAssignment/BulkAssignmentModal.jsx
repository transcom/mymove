import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { Button, Checkbox } from '@trussworks/react-uswds';
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
  const [selectedUsers, setSelectedUsers] = useState({});
  const [isLoading, setIsLoading] = useState(true);
  const [bulkAssignmentData, setBulkAssignmentData] = useState(null);
  const [isSaveDisabled, setIsSaveDisabled] = useState(false);

  const handleCheckboxChange = (userId) => {
    setSelectedUsers((prev) => ({
      ...prev,
      [userId]: !prev[userId],
    }));
  };

  const initUserData = (availableOfficeUsers) => {
    const officeUsers = [];
    const selectedOfficeUsers = {};
    availableOfficeUsers.forEach((user) => {
      const newUserAssignment = {
        ID: user.officeUserId,
        moveAssignments: 0,
      };
      officeUsers.push(newUserAssignment);
      selectedOfficeUsers[user.officeUserId] = true;
    });
    setSelectedUsers(() => selectedOfficeUsers);
    initialValues.userData = officeUsers;
    setIsLoading(false);
  };

  useEffect(() => {
    const fetchData = async () => {
      try {
        getBulkAssignmentData(queueType).then((data) => {
          setBulkAssignmentData(data);
          initUserData(data?.availableOfficeUsers);
          if (data.bulkAssignmentMoveIDs === undefined) {
            setIsSaveDisabled(true);
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

  if (isLoading) return null;

  return (
    <Modal>
      <ModalClose handleClick={() => onClose()} />
      <ModalTitle>
        <h3>
          {title} (
          {bulkAssignmentData?.bulkAssignmentMoveIDs == null ? 0 : bulkAssignmentData?.bulkAssignmentMoveIDs?.length})
        </h3>
      </ModalTitle>
      <div className={styles.BulkAssignmentTable}>
        <Formik
          onSubmit={(values) => {
            const bulkAssignmentSavePayload = values;
            onSubmit({ bulkAssignmentSavePayload });
            onClose();
          }}
          validationSchema={validationSchema}
          initialValues={initialValues}
        >
          {({ handleChange, setValues, values }) => {
            const handleAssignClick = () => {
              const totalMoves = bulkAssignmentData?.bulkAssignmentMoveIDs?.length;
              const numUsers = Object.keys(selectedUsers).filter((id) => selectedUsers[id]).length;
              const baseAssignments = Math.floor(totalMoves / numUsers);
              let remainingMoves = totalMoves % numUsers;

              const newValues = { ...values };

              values.userData.forEach((officeUser) => {
                if (selectedUsers[officeUser.ID]) {
                  const moveAssignments = baseAssignments + (remainingMoves > 0 ? 1 : 0);
                  remainingMoves = Math.max(remainingMoves - 1, 0);
                  newValues.userData.find((u) => u.ID === officeUser.ID).moveAssignments = moveAssignments;
                } else {
                  newValues.userData.find((u) => u.ID === officeUser.ID).moveAssignments = 0;
                }
              });

              setValues({
                ...values,
                ...newValues,
              });
            };
            return (
              <Form>
                <table>
                  <tr>
                    <th>{/* <button type="button">Select All</button> */}</th>
                    <th>User</th>
                    <th>Workload</th>
                    <th>Assignment</th>
                  </tr>
                  {bulkAssignmentData?.availableOfficeUsers?.map((user, i) => {
                    return (
                      <tr key={user.officeUserId}>
                        <td>
                          <input
                            type="checkbox"
                            checked={!!selectedUsers[user.officeUserId]}
                            onChange={() => handleCheckboxChange(user.officeUserId)}
                          />
                        </td>
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
                            name={`userData.${i}.moveAssignments`}
                            id={user.officeUserId}
                            data-testid="assignment"
                            min={0}
                            value={values.userData[i]?.moveAssignments || 0}
                            onChange={(event) => {
                              handleChange(event);

                              const newUserAssignment = {
                                ID: user.officeUserId,
                                moveAssignments: event.target.value ? +event.target.value : 0,
                              };

                              const newUserData = [...values.userData];
                              newUserData[i] = newUserAssignment;

                              setValues({
                                ...values,
                                userData: newUserData,
                              });
                            }}
                          />
                        </td>
                      </tr>
                    );
                  })}
                </table>
                <ModalActions autofocus="true">
                  <div className={styles.BulkAssignmentButtonsContainer}>
                    <div>
                      <Button
                        data-focus="true"
                        className="usa-button--submit"
                        type="submit"
                        data-testid="modalSubmitButton"
                        disabled={isSaveDisabled}
                      >
                        {submitText}
                      </Button>
                      <button
                        className={styles.backbutton}
                        type="button"
                        onClick={() => onClose()}
                        data-testid="modalBackButton"
                      >
                        {closeText}
                      </button>
                    </div>
                    <div>
                      <Button
                        onClick={handleAssignClick}
                        type="button"
                        disabled={!Object.values(selectedUsers).some(Boolean)}
                      >
                        {/* <Button onClick={handleAssignClick} type="button" disabled={isEqualAssignDisabled}> */}
                        Equal Assign
                      </Button>
                    </div>
                  </div>
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

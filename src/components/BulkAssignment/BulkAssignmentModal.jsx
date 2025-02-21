import React, { useCallback, useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import * as Yup from 'yup';

import styles from './BulkAssignmentModal.module.scss';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';
import { getBulkAssignmentData } from 'services/ghcApi';
import { milmoveLogger } from 'utils/milmoveLog';
import { userName } from 'utils/formatters';
import { Form } from 'components/form';

const initialValues = {
  userData: [],
  moveData: [],
};

export const BulkAssignmentModal = ({ onClose, onSubmit, title, submitText, closeText, queueType }) => {
  const [isError, setIsError] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [isDisabled, setIsDisabled] = useState(false);
  const [bulkAssignmentData, setBulkAssignmentData] = useState(null);
  const [selectedUsers, setSelectedUsers] = useState({});
  const [numberOfMoves, setNumberOfMoves] = useState(0);
  const [showCancelModal, setShowCancelModal] = useState(false);

  const errorMessage = 'Cannot assign more moves than are available.';

  const handleCheckboxChange = (userId) => {
    setSelectedUsers((prev) => ({
      ...prev,
      [userId]: !prev[userId],
    }));
  };

  const isAllSelected = () => {
    return Object.keys(selectedUsers).every((id) => selectedUsers[id]);
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

  const fetchData = useCallback(async () => {
    try {
      const data = await getBulkAssignmentData(queueType);
      setBulkAssignmentData(data);
      initUserData(data?.availableOfficeUsers);

      if (!data.bulkAssignmentMoveIDs) {
        setIsDisabled(true);
        setNumberOfMoves(0);
      } else {
        setNumberOfMoves(data.bulkAssignmentMoveIDs.length);
      }
    } catch (err) {
      milmoveLogger.error('Error fetching bulk assignment data:', err);
    }
  }, [queueType]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  initialValues.moveData = bulkAssignmentData?.bulkAssignmentMoveIDs;

  const validationSchema = Yup.object().shape({
    assignment: Yup.number().min(0).typeError('Assignment must be a number'),
  });

  if (isLoading) return null;

  return (
    <div>
      <Modal className={styles.BulkModal}>
        {!showCancelModal && <ModalClose handleClick={() => setShowCancelModal(true)} />}
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

              const bulkAssignmentSavePayload = {
                moveData: values.moveData,
                userData: values.userData.filter((user) => user.moveAssignments > 0),
              };
              onSubmit({ bulkAssignmentSavePayload });
              onClose();
            }}
            validationSchema={validationSchema}
            initialValues={initialValues}
          >
            {({ handleChange, setValues, values }) => {
              const handleAssignmentChange = (event, user, i) => {
                handleChange(event);
                setIsError(false);

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
              };

              const handleEqualAssignClick = () => {
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
                      <th>
                        <input
                          data-testId="selectDeselectAllButton"
                          type="checkbox"
                          checked={isAllSelected()}
                          onChange={() => {
                            const allSelected = Object.keys(selectedUsers).every((id) => selectedUsers[id]);
                            const newSelectedUsers = {};

                            bulkAssignmentData.availableOfficeUsers.forEach((user) => {
                              newSelectedUsers[user.officeUserId] = !allSelected;
                            });

                            setSelectedUsers(newSelectedUsers);
                          }}
                        />
                      </th>
                      <th className={styles.UserNameHeader}>Users</th>
                      <th>Workload</th>
                      <th>Assignment</th>
                    </tr>
                    {bulkAssignmentData?.availableOfficeUsers?.map((user, i) => {
                      return (
                        <tr key={user}>
                          <td>
                            <input
                              data-testid="bulkAssignmentUserCheckbox"
                              type="checkbox"
                              checked={!!selectedUsers[user.officeUserId]}
                              onChange={() => handleCheckboxChange(user.officeUserId)}
                            />
                          </td>
                          <td>
                            <p data-testid="bulkAssignmentUser" className={styles.officeUserFormattedName}>
                              {userName(user)}
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
                              onChange={(event) => handleAssignmentChange(event, user, i)}
                            />
                          </td>
                        </tr>
                      );
                    })}
                  </table>
                  {showCancelModal ? (
                    <div className={styles.areYouSureSection}>
                      <small className={styles.hint}>
                        Any unsaved work will be lost. Are you sure you want to cancel?
                      </small>
                      <div className={styles.confirmButtons}>
                        <Button
                          className={styles.cancelNoButton}
                          data-testid="cancelModalNo"
                          onClick={() => setShowCancelModal(false)}
                        >
                          No
                        </Button>
                        <Button
                          className={styles.cancelYesButton}
                          data-testid="cancelModalYes"
                          secondary
                          onClick={onClose}
                        >
                          Discard Changes
                        </Button>
                      </div>
                    </div>
                  ) : (
                    <ModalActions autofocus="true">
                      <div className={styles.BulkAssignmentButtonsContainer}>
                        <div>
                          <Button disabled={isDisabled} data-focus="true" type="submit" data-testid="modalSubmitButton">
                            {submitText}
                          </Button>
                          <Button
                            type="button"
                            className={styles.button}
                            unstyled
                            onClick={() => setShowCancelModal(true)}
                            data-testid="modalCancelButton"
                          >
                            {closeText}
                          </Button>
                        </div>
                        <div>
                          <Button
                            onClick={handleEqualAssignClick}
                            type="button"
                            data-testid="modalEqualAssignButton"
                            disabled={!Object.values(selectedUsers).some(Boolean)}
                          >
                            Equal Assign
                          </Button>
                        </div>
                      </div>
                      {isError && <div className={styles.errorMessage}>{errorMessage}</div>}
                    </ModalActions>
                  )}
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

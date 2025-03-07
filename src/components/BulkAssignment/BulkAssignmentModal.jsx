import React, { useCallback, useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { Alert, Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import { Switch } from '@material-ui/core';
import * as Yup from 'yup';

import styles from './BulkAssignmentModal.module.scss';

import { isBooleanFlagEnabled } from 'utils/featureFlags';
import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';
import { getBulkAssignmentData } from 'services/ghcApi';
import { milmoveLogger } from 'utils/milmoveLog';
import { userName } from 'utils/formatters';
import { Form } from 'components/form';

const initialValues = {
  userData: [],
  moveData: [],
};

export const BulkAssignmentModal = ({ onClose, onSubmit, submitText, closeText, queueType }) => {
  const bulkAssignmentSwitchLabels = ['Bulk Assignment', 'Bulk Re-assignment'];

  const [isError, setIsError] = useState(false);
  const [isBulkReAssignmentMode, setIsBulkReAssignmentMode] = useState(false);
  const [bulkAssignmentData, setBulkAssignmentData] = useState(null);
  const [isDisabled, setIsDisabled] = useState(false);
  const [showCancelModal, setShowCancelModal] = useState(false);
  const [numberOfMoves, setNumberOfMoves] = useState(0);
  const [selectedUsers, setSelectedUsers] = useState({});
  const [selectedRadio, setSelectedRadio] = useState(null);
  const errorMessage = 'Cannot assign more moves than are available.';
  const [isBulkReAssignmentEnabled, setIsBulkReAssignmentEnabled] = useState(false);

  useEffect(() => {
    // checking feature flag to see if DODID input should be disabled
    // this data pulls from Okta and doens't let the customer update it
    const fetchData = async () => {
      setIsBulkReAssignmentEnabled(await isBooleanFlagEnabled('bulk_re_assignment'));
    };
    fetchData();
  }, []);

  const handleRadioChange = (index) => {
    setSelectedRadio(index);

    setSelectedUsers((prev) => ({
      ...prev,
      [index]: false,
    }));

    if (isBulkReAssignmentMode) {
      const reAssignableMoves = bulkAssignmentData.availableOfficeUsers.find(
        (user) => user.officeUserId === index,
      ).workload;
      setNumberOfMoves(reAssignableMoves);
    }
  };

  const handleCheckboxChange = (userId) => {
    setSelectedUsers((prev) => ({
      ...prev,
      [userId]: !prev[userId],
    }));
  };

  const handleAssignmentModeChange = (event) => {
    setIsBulkReAssignmentMode(event.target.checked);
    if (event.target.checked && selectedRadio != null) {
      const reAssignableMoves = bulkAssignmentData.availableOfficeUsers.find(
        (user) => user.officeUserId === selectedRadio,
      ).workload;
      setNumberOfMoves(reAssignableMoves);
    } else if (event.target.checked && selectedRadio == null) {
      setNumberOfMoves(0);
    } else {
      setNumberOfMoves(bulkAssignmentData.bulkAssignmentMoveIDs.length);
    }

    if (!event.target.checked) {
      setSelectedRadio(null);
    }

    // need to sync workload with user load
    const newValues = { ...initialValues };
    initialValues.userData.forEach((element) => {
      if (event.target.checked) {
        const userWorkload = bulkAssignmentData.availableOfficeUsers.find(
          (user) => user.officeUserId === element.ID,
        ).workload;
        newValues.userData.find((u) => u.ID === element.ID).moveAssignments = userWorkload;
      } else {
        newValues.userData.find((u) => u.ID === element.ID).moveAssignments = 0;
      }
    });
    initialValues.userData = newValues.userData;
  };

  const isAllSelected = () => {
    const selectedIds = Object.keys(selectedUsers);
    return selectedIds.length > 0 && selectedIds.every((id) => selectedUsers[id]);
  };

  const isFormUnchanged = (values) => {
    return values.userData.every((user) => user.moveAssignments === 0);
  };

  const handleCancelClick = (values) => () => {
    setIsError(false);
    if (isFormUnchanged(values)) onClose();
    else setShowCancelModal(true);
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
        setIsDisabled(false);
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

  return (
    <div>
      <Modal className={styles.BulkModal}>
        <div className={styles.BulkAssignmentTable}>
          <Formik
            onSubmit={(values) => {
              const totalAssignment = values?.userData?.reduce((sum, item) => sum + item.moveAssignments, 0);

              const totalAssignedMovesGreaterThanMovesAvailableReassignment = totalAssignment > numberOfMoves;
              if (totalAssignedMovesGreaterThanMovesAvailableReassignment) {
                setIsError(true);
                return;
              }
              if (totalAssignment === 0) {
                onClose();
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
                <>
                  {!showCancelModal && <ModalClose handleClick={handleCancelClick(values)} />}
                  <ModalTitle>
                    <h3>
                      {isBulkReAssignmentMode ? bulkAssignmentSwitchLabels[1] : bulkAssignmentSwitchLabels[0]} (
                      {numberOfMoves})
                    </h3>
                  </ModalTitle>
                  {isBulkReAssignmentEnabled && (
                    <Switch
                      data-testid="modalReAssignModeToggleSwitch"
                      name="BulkAssignmentModeSwitch"
                      onChange={handleAssignmentModeChange}
                      inputProps={{ 'aria-label': 'BulkAssignmentModeSwitch' }}
                    />
                  )}
                  <Form>
                    <table>
                      <tr>
                        <th>
                          {!isBulkReAssignmentEnabled && (
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
                          )}
                        </th>
                        <th className={styles.UserNameHeader}>User</th>
                        <th>Workload</th>
                        <th>Assignment</th>
                        {isBulkReAssignmentMode && <th>Re-assign User</th>}
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
                                disabled={selectedRadio === user.officeUserId}
                                onChange={(event) => handleAssignmentChange(event, user, i)}
                              />
                            </td>
                            {isBulkReAssignmentMode && (
                              <td className={styles.BulkAssignmentDataCenter}>
                                <input
                                  type="radio"
                                  name={`userData.${i}.moveReAssignment`}
                                  value={user.officeUserId}
                                  checked={selectedRadio === user.officeUserId}
                                  onChange={() => handleRadioChange(user.officeUserId)}
                                />
                              </td>
                            )}
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
                          <div className={styles.BulkAssignmentButtonsLeft}>
                            <Button
                              disabled={isDisabled || isFormUnchanged(values)}
                              data-focus="true"
                              type="submit"
                              data-testid="modalSubmitButton"
                            >
                              {submitText}
                            </Button>
                            <Button
                              type="button"
                              className={styles.button}
                              unstyled
                              onClick={handleCancelClick(values)}
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
                              disabled={!Object.values(selectedUsers).some(Boolean) || isDisabled}
                            >
                              Equal Assign
                            </Button>
                          </div>
                        </div>
                      </ModalActions>
                    )}
                    {isError && (
                      <Alert type="error" headingLevel="h4" slim>
                        {errorMessage}
                      </Alert>
                    )}
                  </Form>
                </>
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

  submitText: PropTypes.string,
  closeText: PropTypes.string,
};

BulkAssignmentModal.defaultProps = {
  submitText: 'Save',
  closeText: 'Cancel',
};

BulkAssignmentModal.displayName = 'BulkAssignmentModal';

export default connectModal(BulkAssignmentModal);

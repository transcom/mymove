import React, { useCallback, useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { Alert, Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import { Switch, FormControlLabel } from '@material-ui/core';
import * as Yup from 'yup';

import styles from './BulkAssignmentModal.module.scss';

import { FEATURE_FLAG_KEYS } from 'shared/constants';
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

export const BulkAssignmentModal = ({ onClose, onSubmit, submitText, closeText, queueType, activeOfficeID }) => {
  const bulkAssignmentSwitchLabels = ['Bulk Assignment', 'Bulk Re-assignment'];
  const errorMessage = 'Cannot assign more moves than are available.';

  const [isError, setIsError] = useState(false);
  const [isBulkReAssignmentMode, setIsBulkReAssignmentMode] = useState(false);
  const [bulkAssignmentData, setBulkAssignmentData] = useState(null);
  const [isDisabled, setIsDisabled] = useState(false);
  const [showCancelModal, setShowCancelModal] = useState(false);
  const [numberOfMoves, setNumberOfMoves] = useState(0);
  const [selectedUsers, setSelectedUsers] = useState({});
  const [selectedRadio, setSelectedRadio] = useState(null);
  const [isBulkReAssignmentEnabled, setIsBulkReAssignmentEnabled] = useState(false);

  useEffect(() => {
    const fetchFlag = async () => {
      setIsBulkReAssignmentEnabled(await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.BULK_RE_ASSIGNMENT));
    };
    fetchFlag();
  }, []);

  const handleCheckboxChange = (userId) => {
    setSelectedUsers((prev) => ({
      ...prev,
      [userId]: !prev[userId],
    }));
  };

  const isAllSelected = () => {
    return Object.keys(selectedUsers).every((id) => selectedUsers[id]);
  };

  const isFormUnchanged = (values) => {
    return values.userData.every((user) => user.moveAssignments === 0);
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
    setSelectedUsers(selectedOfficeUsers);
    initialValues.userData = officeUsers;
  };

  const fetchData = useCallback(async () => {
    try {
      const data = await getBulkAssignmentData(queueType, activeOfficeID);
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
  }, [queueType, activeOfficeID]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  initialValues.moveData = bulkAssignmentData?.bulkAssignmentMoveIDs;

  const validationSchema = Yup.object().shape({
    assignment: Yup.number().min(0).typeError('Assignment must be a number'),
  });

  // helper to set all checkboxes to a given value
  const setAllCheckboxes = (value) => {
    if (bulkAssignmentData?.availableOfficeUsers) {
      const newSelectedUsers = {};
      bulkAssignmentData.availableOfficeUsers.forEach((user) => {
        newSelectedUsers[user.officeUserId] = value;
      });
      setSelectedUsers(newSelectedUsers);
    }
  };

  // define cancel click behavior
  const handleCancelClick = (values) => () => {
    if (!isFormUnchanged(values)) {
      setShowCancelModal(true);
    } else {
      onClose();
    }
  };

  return (
    <div>
      <Modal className={styles.BulkModal} onClose={onClose}>
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
              const handleRadioChange = (index) => {
                // to avoid confusion between current and previous selections
                const newSelection = index;
                setSelectedRadio(newSelection);

                if (isBulkReAssignmentMode) {
                  const reAssignableMoves = bulkAssignmentData.availableOfficeUsers.find(
                    (user) => user.officeUserId === newSelection,
                  ).workload;
                  setNumberOfMoves(reAssignableMoves);

                  // need to reset assignment entries between re-assignment changes
                  const newValues = { ...initialValues };
                  newValues.userData.forEach((element) => {
                    newValues.userData.find((u) => u.ID === element.ID).moveAssignments = 0;
                  });
                  setValues({
                    ...values,
                    ...newValues,
                  });
                }
              };
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

              const handleAssignmentModeChange = (event) => {
                setIsBulkReAssignmentMode(event.target.checked);
                if ((event.target.checked && selectedRadio == null) || !bulkAssignmentData?.bulkAssignmentMoveIDs) {
                  setNumberOfMoves(0);
                  setAllCheckboxes(true);
                } else {
                  setNumberOfMoves(bulkAssignmentData.bulkAssignmentMoveIDs.length);
                  setAllCheckboxes(true);
                }

                if (!event.target.checked) {
                  setSelectedRadio(null);
                }

                // reset assignment entries between form mode changes
                const newValues = { ...initialValues };
                newValues.userData.forEach((element) => {
                  const userEntry = newValues.userData.find((u) => u.ID === element.ID);
                  if (userEntry) {
                    userEntry.moveAssignments = 0;
                  }
                });
                setValues({
                  ...values,
                  ...newValues,
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
                    <h3 data-testid="modalTitleHeader">
                      {isBulkReAssignmentMode ? bulkAssignmentSwitchLabels[1] : bulkAssignmentSwitchLabels[0]} (
                      {numberOfMoves})
                    </h3>
                  </ModalTitle>
                  {isBulkReAssignmentEnabled && (
                    <FormControlLabel
                      control={
                        <Switch
                          data-testid="modalReAssignModeToggleSwitch"
                          name="BulkAssignmentModeSwitch"
                          onChange={handleAssignmentModeChange}
                          inputProps={{ 'aria-label': 'BulkAssignmentModeSwitch' }}
                          color="primary"
                        />
                      }
                      label="Re-assignment"
                    />
                  )}
                  <Form>
                    <table>
                      <thead>
                        <tr>
                          <th>
                            <input
                              data-testid="selectDeselectAllButton"
                              hidden={isBulkReAssignmentMode}
                              type="checkbox"
                              checked={isAllSelected()}
                              onChange={() => {
                                const allSelected = isAllSelected();
                                if (bulkAssignmentData?.availableOfficeUsers) {
                                  const newSelectedUsers = {};
                                  bulkAssignmentData.availableOfficeUsers.forEach((user) => {
                                    newSelectedUsers[user.officeUserId] = !allSelected;
                                  });
                                  setSelectedUsers(newSelectedUsers);
                                }
                              }}
                            />
                          </th>
                          <th>User</th>
                          <th>Current Workload</th>
                          <th>Assignment</th>
                          {isBulkReAssignmentMode && <th>Re-assign Workload</th>}
                        </tr>
                      </thead>
                      <tbody>
                        {bulkAssignmentData?.availableOfficeUsers?.map((user, i) => (
                          <tr key={user.officeUserId}>
                            <td>
                              <input
                                data-testid="bulkAssignmentUserCheckbox"
                                hidden={isBulkReAssignmentMode}
                                type="checkbox"
                                checked={selectedUsers[user.officeUserId] && selectedRadio !== user.officeUserId}
                                disabled={
                                  selectedRadio === user.officeUserId ||
                                  (isBulkReAssignmentMode && selectedRadio == null)
                                }
                                onChange={() => handleCheckboxChange(user.officeUserId)}
                              />
                            </td>
                            <td>
                              <p data-testid="bulkAssignmentUser">{userName(user)}</p>
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
                                value={values.userData[i]?.moveAssignments.toString() || 0}
                                disabled={
                                  selectedRadio === user.officeUserId ||
                                  (isBulkReAssignmentMode && selectedRadio == null)
                                }
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
                        ))}
                      </tbody>
                    </table>
                    {showCancelModal ? (
                      <div className={styles.areYouSureSection}>
                        <small className={styles.hint}>
                          Any unsaved work will be lost. Are you sure you want to cancel?
                        </small>
                        <div className={styles.confirmButtons}>
                          <Button outline data-testid="cancelModalNo" onClick={() => setShowCancelModal(false)}>
                            No
                          </Button>
                          <Button
                            className="usa-button-destructive"
                            secondary
                            data-testid="cancelModalYes"
                            onClick={onClose}
                          >
                            Discard Changes
                          </Button>
                        </div>
                      </div>
                    ) : (
                      <ModalActions>
                        <div className={styles.BulkAssignmentButtonsContainer}>
                          <div>
                            <Button
                              onClick={handleEqualAssignClick}
                              type="button"
                              data-testid="modalEqualAssignButton"
                              outline
                              hidden={isBulkReAssignmentMode}
                              disabled={!Object.values(selectedUsers).some(Boolean)}
                            >
                              Equal Assign
                            </Button>
                          </div>
                          <div className={styles.BulkAssignmentButtonsRight}>
                            <Button
                              type="button"
                              outline
                              onClick={handleCancelClick(values)}
                              data-testid="modalCancelButton"
                            >
                              {closeText}
                            </Button>
                            <Button
                              disabled={isDisabled || isFormUnchanged(values)}
                              data-focus="true"
                              type="submit"
                              data-testid="modalSubmitButton"
                            >
                              {submitText}
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

import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { Button, Checkbox } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import { Switch, FormGroup, FormControlLabel, FormLabel, FormControl } from '@material-ui/core';
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
  const bulkAssignmentSwitchLabels = ['Bulk Assignment', 'Bulk Re-assignment'];

  const [isLoading, setIsLoading] = useState(true);
  const [isError, setIsError] = useState(false);
  const [isBulkReAssignmentMode, setIsBulkReAssignmentMode] = useState(false);
  const [bulkAssignmentData, setBulkAssignmentData] = useState(null);
  const [isSaveDisabled, setIsSaveDisabled] = useState(false);
  const [numberOfMoves, setNumberOfMoves] = useState(0);
  const [selectedUsers, setSelectedUsers] = useState({});
  const [selectedRadio, setSelectedRadio] = useState(null);

  const errorMessage = 'Cannot assign more moves than are available.';

  const handleRadioChange = (index) => {
    setSelectedRadio(index);
  };

  const handleCheckboxChange = (userId) => {
    setSelectedUsers((prev) => ({
      ...prev,
      [userId]: !prev[userId],
    }));
  };

  const handleAssignmentModeChange = (event) => {
    setIsBulkReAssignmentMode(event.target.checked);
  };
  const isAllSelected = () => {
    const selectedIds = Object.keys(selectedUsers);
    return selectedIds.length > 0 && selectedIds.every((id) => selectedUsers[id]);
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

  if (isLoading) return null;

  return (
    <Modal>
      <ModalClose handleClick={() => onClose()} />
      <ModalTitle>
        <h3>
          {isBulkReAssignmentMode ? bulkAssignmentSwitchLabels[1] : bulkAssignmentSwitchLabels[0]} ({numberOfMoves})
        </h3>
      </ModalTitle>
      <Switch name="BulkAssignmentModeSwitch" onChange={handleAssignmentModeChange} />
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
            // const handleAssignmentChange = (event, i) => {
            //   handleChange(event);
            //   setIsError(false);

            //   let newUserAssignment;
            //   if (event.target.value !== '') {
            //     newUserAssignment = {
            //       ID: event.target.id,
            //       moveAssignments: +event.target.value,
            //     };
            //   } else {
            //     newUserAssignment = {
            //       ID: event.target.id,
            //       moveAssignments: 0,
            //     };
            //   }

            //   const newValues = values;
            //   newValues.userData[i] = newUserAssignment;

            //   setValues({
            //     ...values,
            //     userData: newValues.userData,
            //   });
            // };

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

            const handleReAssignClick = () => {
              const totalMoves = bulkAssignmentData?.bulkAssignmentMoveIDs?.length;
              const numUsers = Object.keys(selectedUsers).filter((id) => selectedUsers[id]).length;
              const baseAssignments = Math.floor(totalMoves / numUsers);

              let remainingMoves = totalMoves % numUsers;

              const newValues = { ...values };

              values.userData.forEach((officeUser) => {
                if (selectedUsers[officeUser.ID] && selectedUsers[officeUser.ID] !== selectedRadio) {
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
                    <th>User</th>
                    <th>Workload</th>
                    <th>Assignment</th>
                    {isBulkReAssignmentMode && <th>Re-assign User</th>}
                  </tr>
                  {bulkAssignmentData?.availableOfficeUsers?.map((user, i) => {
                    return (
                      <tr key={user.officeUserId}>
                        <td>
                          <input
                            data-testid="bulkAssignmentUserCheckbox"
                            type="checkbox"
                            checked={!!selectedUsers[user.officeUserId] && selectedRadio !== user.officeUserId}
                            disabled={selectedRadio === user.officeUserId}
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
                        {isBulkReAssignmentMode && (
                          <td className={styles.BulkAssignmentDataCenter}>
                            <input
                              type="radio"
                              name="chooseReAssignment"
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
                      {isBulkReAssignmentMode ? (
                        <Button
                          onClick={handleReAssignClick}
                          type="button"
                          data-testid="modalEqualAssignButton"
                          disabled={!Object.values(selectedUsers).some(Boolean) || selectedRadio === null}
                        >
                          Equal Re-Assign
                        </Button>
                      ) : (
                        <Button
                          onClick={handleAssignClick}
                          type="button"
                          data-testid="modalEqualAssignButton"
                          disabled={!Object.values(selectedUsers).some(Boolean)}
                        >
                          Equal Assign
                        </Button>
                      )}
                    </div>
                  </div>
                </ModalActions>
                {isError && <div className={styles.errorMessage}>{errorMessage}</div>}
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

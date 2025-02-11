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
  const [bulkAssignmentData, setBulkAssignmentData] = useState(null);
  const [isDisabled, setIsDisabled] = useState(false);
  const [numberOfMoves, setNumberOfMoves] = useState(0);
  const fetchData = useCallback(async () => {
    try {
      const data = await getBulkAssignmentData(queueType);
      setBulkAssignmentData(data);

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
              const handleAssignmentChange = (event, i) => {
                handleChange(event);

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
                        <tr key={user}>
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
                      onClick={onClose}
                      data-testid="modalCancelButton"
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

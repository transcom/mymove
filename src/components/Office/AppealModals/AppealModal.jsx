import React from 'react';
import { Formik, Field } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button, Fieldset, Grid, Radio } from '@trussworks/react-uswds';

import styles from './AppealModal.module.scss';

import { Form } from 'components/form';
import formStyles from 'styles/form.module.scss';
import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';
import TextField from 'components/form/fields/TextField/TextField';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Hint from 'components/Hint';

export const AppealModal = ({ onClose, onSubmit, selectedReportViolation, isSeriousIncidentAppeal }) => {
  const appealSchema = Yup.object().shape({
    remarks: Yup.string().required('Remarks are required'),
    appealStatus: Yup.string().required('Appeal status is required'),
  });

  const initialValues = {
    remarks: '',
    appealStatus: '',
  };

  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.AppealModal}>
          <ModalClose handleClick={onClose} />
          <ModalTitle>
            <h2 className={styles.ModalTitle} data-testid="appealModalTitle">
              Leave Appeal Decision
            </h2>
            {isSeriousIncidentAppeal ? (
              <Hint data-testid="seriousIncidentModalHint">Serious Incident</Hint>
            ) : (
              <Hint data-testid="violationModalHint">
                {selectedReportViolation?.violation?.paragraphNumber} {selectedReportViolation?.violation?.title}
              </Hint>
            )}
          </ModalTitle>
          <Formik validationSchema={appealSchema} onSubmit={onSubmit} initialValues={initialValues} validateOnMount>
            {({ isValid }) => {
              return (
                <Form className={formStyles.form}>
                  <Fieldset>
                    <Grid row>
                      <Grid col={12}>
                        <TextField
                          label="Remarks"
                          id="remarks"
                          name="remarks"
                          display="textarea"
                          data-testid="addAppealRemarks"
                        />
                      </Grid>
                    </Grid>
                  </Fieldset>
                  <Fieldset>
                    <Field
                      as={Radio}
                      id="sustainedRadio"
                      label="Sustained"
                      name="appealStatus"
                      value="sustained"
                      data-testid="sustainedRadio"
                    />
                    <Field
                      as={Radio}
                      id="rejectedRadio"
                      label="Rejected"
                      name="appealStatus"
                      value="rejected"
                      data-testid="rejectedRadio"
                    />
                  </Fieldset>
                  <ModalActions>
                    <Button type="submit" disabled={!isValid}>
                      Save
                    </Button>
                    <Button
                      type="button"
                      onClick={() => onClose()}
                      data-testid="modalCancelButton"
                      secondary
                      className={styles.CancelButton}
                    >
                      Cancel
                    </Button>
                  </ModalActions>
                </Form>
              );
            }}
          </Formik>
        </Modal>
      </ModalContainer>
    </div>
  );
};

AppealModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

AppealModal.displayName = 'AppealModal';

export default connectModal(AppealModal);

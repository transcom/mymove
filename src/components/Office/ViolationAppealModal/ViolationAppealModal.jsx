import React from 'react';
import { Formik, Field } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button, Fieldset, Grid, Radio } from '@trussworks/react-uswds';

import styles from './ViolationAppealModal.module.scss';

import { Form } from 'components/form';
import formStyles from 'styles/form.module.scss';
import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';
import TextField from 'components/form/fields/TextField/TextField';

export const ViolationAppealModal = ({ onClose, onSubmit }) => {
  const violationAppealSchema = Yup.object().shape({
    remarks: Yup.string().required('Remarks are required'),
    appealStatus: Yup.string().required('Appeal status is required'),
  });

  const initialValues = {
    remarks: '',
    appealStatus: '',
  };

  return (
    <Modal className={styles.EditFacilityInfoModal}>
      <ModalClose handleClick={onClose} />
      <ModalTitle>
        <h2 className={styles.ModalTitle}>Add Violation Appeal</h2>
      </ModalTitle>
      <Formik
        validationSchema={violationAppealSchema}
        onSubmit={onSubmit}
        initialValues={initialValues}
        validateOnMount
      >
        {({ isValid }) => {
          return (
            <Form className={formStyles.form}>
              <Fieldset>
                <Grid row>
                  <Grid col={12}>
                    <TextField label="Remarks" id="remarks" name="remarks" />
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
  );
};

ViolationAppealModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

export default connectModal(ViolationAppealModal);

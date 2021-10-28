import React from 'react';
import { Formik, Field } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button, Label, Textarea } from '@trussworks/react-uswds';

import styles from './FinancialReviewModal.module.scss';

import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';

const reviewSITExtensionSchema = Yup.object().shape({
  requestReason: Yup.string().required('Required'),
  daysApproved: Yup.number()
    .min(1, 'Additional days approved must be greater than or equal to 1.')
    .required('Required'),
  officeRemarks: Yup.string().nullable(),
});

const FinancialReviewModal = ({ onClose, onSubmit, summarySITComponent }) => {
  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.FinancialReviewModal}>
          <ModalClose handleClick={() => onClose()} />
          <ModalTitle>
            <h2>Edit SIT authorization</h2>
          </ModalTitle>
          <div className={styles.summarySITComponent}>{summarySITComponent}</div>
          <div className={styles.ModalPanel}>
            <Formik
              validationSchema={reviewSITExtensionSchema}
              onSubmit={(e) => onSubmit(e)}
              initialValues={{
                requestReason: '',
                daysApproved: '',
                officeRemarks: '',
              }}
            >
              {({ isValid }) => {
                return (
                  <Form>
                    <Label htmlFor="remarks">Remarks</Label>
                    <Field as={Textarea} data-testid="remarks" label="No" name="remarks" id="remarks" />
                    <ModalActions>
                      <Button type="submit" disabled={!isValid}>
                        Save
                      </Button>
                      <Button
                        type="button"
                        onClick={() => onClose()}
                        data-testid="modalCancelButton"
                        outline
                        className={styles.CancelButton}
                      >
                        Cancel
                      </Button>
                    </ModalActions>
                  </Form>
                );
              }}
            </Formik>
          </div>
        </Modal>
      </ModalContainer>
    </div>
  );
};

FinancialReviewModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  summarySITComponent: PropTypes.node.isRequired,
};
export default FinancialReviewModal;

import React, { useState } from 'react';
import classnames from 'classnames';
import { Formik, Field } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button, Label, Textarea, Checkbox } from '@trussworks/react-uswds';

import styles from './FinancialReviewModal.module.scss';

import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';

const financialReviewSchema = Yup.object().shape({
  remarks: Yup.string().required('Required'),
  reviewCheckbox: Yup.boolean().oneOf([true], 'Must click needs review checkbox'),
});

const FinancialReviewModal = ({ onClose, onSubmit }) => {
  const [remarksDisabled, setremarksDisabled] = useState(true);
  const labelClass = classnames({
    [styles.RemarksLabelDisabled]: remarksDisabled,
  });
  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.FinancialReviewModal}>
          <ModalClose handleClick={() => onClose()} />
          <ModalTitle>
            <h2>Flag for Financial Review</h2>
          </ModalTitle>
          <p>This will let the financial office know to review this move for potential costs to the customer.</p>
          <div>
            <Formik
              validationSchema={financialReviewSchema}
              remarksDisabled
              onSubmit={(e) => onSubmit(e)}
              initialValues={{
                remarks: '',
                reviewCheckbox: false,
              }}
            >
              {({ isValid }) => {
                return (
                  <Form>
                    <Checkbox
                      data-testid="reviewCheckbox"
                      label="This move needs financial review"
                      name="reviewCheckbox"
                      onChange={() => {
                        setremarksDisabled(!remarksDisabled);
                      }}
                      id="reviewCheckbox"
                    />
                    <Label className={labelClass} htmlFor="remarks">
                      Remarks
                    </Label>
                    <Field
                      disabled={remarksDisabled}
                      as={Textarea}
                      data-testid="remarks"
                      label="No"
                      name="remarks"
                      id="remarks"
                      className={styles.RemarksField}
                    />
                    <ModalActions>
                      <Button type="submit" disabled={isValid}>
                        Save
                      </Button>
                      <Button
                        type="button"
                        onClick={() => onClose()}
                        data-testid="modalCancelButton"
                        outline
                        className="usa-button--tertiary"
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
};
export default FinancialReviewModal;

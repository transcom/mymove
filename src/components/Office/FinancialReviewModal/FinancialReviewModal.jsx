import React from 'react';
import classnames from 'classnames';
import { Formik, Field } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button, Label, Textarea, Radio, FormGroup } from '@trussworks/react-uswds';

import styles from './FinancialReviewModal.module.scss';

import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';

const financialReviewSchema = Yup.object().shape({
  remarks: Yup.string().test('remarks', 'Remarks are required', (value) => value?.length > 0),
  // must select yest or no before they can click save.
  flagForReview: Yup.string().required('Required').oneOf(['yes']),
});

function FinancialReviewModal({ onClose, onSubmit }) {
  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.FinancialReviewModal}>
          <ModalClose handleClick={onClose} />
          <ModalTitle>
            <h2>Does this move need financial review?</h2>
          </ModalTitle>
          <div>
            <Formik
              initialValues={{
                remarks: '',
                flagForReview: '',
              }}
              validationSchema={financialReviewSchema}
              onSubmit={(values) => onSubmit(values.remarks)}
              validateOnMount
            >
              {({ values, isValid }) => {
                const { flagForReview } = values;
                return (
                  <Form>
                    <FormGroup>
                      <div>Select Yes to flag this move for financial review from the financial review office.</div>
                      <div>Enter remarks to give more detail.</div>
                      <div>
                        <Field
                          as={Radio}
                          label="Yes"
                          id="flagForReview"
                          name="flagForReview"
                          value="yes"
                          title="Yes"
                          type="radio"
                        />
                        <Field
                          as={Radio}
                          label="No"
                          id="doNotFlagforReview"
                          name="flagForReview"
                          title="No"
                          value="no"
                          type="radio"
                        />
                      </div>
                    </FormGroup>
                    <Label
                      className={classnames({
                        [styles.RemarksLabelDisabled]: flagForReview !== 'yes',
                      })}
                      htmlFor="remarks"
                    >
                      Remarks for financial office
                    </Label>
                    {/* Need to set remarks to nothing when no is selected */}
                    <Field
                      disabled={!(flagForReview === 'yes')}
                      as={Textarea}
                      data-testid="remarks"
                      label="No"
                      name="remarks"
                      id="remarks"
                      className={styles.RemarksField}
                    />
                    <ModalActions>
                      <Button type="submit" disabled={!isValid}>
                        Save
                      </Button>
                      <Button type="button" onClick={onClose} outline className="usa-button--tertiary">
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
}

FinancialReviewModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};
export default FinancialReviewModal;

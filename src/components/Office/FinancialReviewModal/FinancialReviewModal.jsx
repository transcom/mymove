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
  remarks: Yup.string().when('flagForReview', {
    is: 'yes',
    then: Yup.string().required('Remarks are required'),
  }),
  flagForReview: Yup.string().required('Required').oneOf(['yes', 'no']),
});

function FinancialReviewModal({ remarks, onClose, onSubmit }) {
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
                remarks,
                flagForReview: 'yes',
              }}
              validationSchema={financialReviewSchema}
              onSubmit={(values) => onSubmit(values.flagForReview, values.remarks)}
              validateOnMount
            >
              {({ values, isValid }) => {
                const { flagForReview } = values;

                return (
                  <Form>
                    <FormGroup>
                      <div>
                        Select <strong>Yes</strong> to flag this move for review from the service&apos;s financial
                        office. Enter remarks to give more detail.
                      </div>
                      <div>
                        <Field
                          as={Radio}
                          label="Yes, flag for financial review"
                          id="flagForReview"
                          name="flagForReview"
                          value="yes"
                          checked={flagForReview === 'yes'}
                        />
                        <Field
                          as={Radio}
                          label="No"
                          id="notFlagForReview"
                          name="flagForReview"
                          title="No, do not flag for financial review"
                          value="no"
                          checked={flagForReview === 'no'}
                        />
                      </div>
                    </FormGroup>
                    <Label
                      className={classnames({
                        [styles.RemarksLabelDisabled]: flagForReview === 'no',
                      })}
                      htmlFor="remarks"
                    >
                      Remarks for financial office
                    </Label>
                    <Field
                      disabled={flagForReview === 'no'}
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
  remarks: PropTypes.string,
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

FinancialReviewModal.defaultProps = {
  remarks: '',
};

export default FinancialReviewModal;

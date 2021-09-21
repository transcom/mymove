import React from 'react';
import classnames from 'classnames';
import { Formik, Field } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button, Radio, FormGroup, Label, Textarea } from '@trussworks/react-uswds';

import { SITExtensionShape } from '../../../types/sitExtensions';

import styles from './ReviewSITExtensionModal.module.scss';

import MaskedTextField from 'components/form/fields/MaskedTextField';
import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
import { sitExtensionReasons } from 'constants/sitExtensions';

const ReviewSITExtensionsModal = ({ onClose, onSubmit, sitExtension }) => {
  const reviewSITExtensionSchema = Yup.object().shape({
    acceptExtension: Yup.mixed().oneOf(['yes', 'no']).required('Required'),
    daysApproved: Yup.number().when('acceptExtension', {
      is: 'yes',
      then: Yup.number()
        .min(1, 'Days approved must be greater than or equal to 1.')
        .max(sitExtension.requestedDays, 'Days approved must be equal to or less than the requested days.')
        .required('Required'),
    }),
    officeRemarks: Yup.string().nullable(),
  });

  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.ReviewSITExtensionModal}>
          <ModalClose handleClick={() => onClose()} />
          <ModalTitle>
            <h2>Review request for extension</h2>
          </ModalTitle>
          <div className={styles.ModalPanel}>
            <div className={styles.SITSummary}>
              <div>
                <dt>Additional days requested:</dt>
                <dd>{sitExtension.requestedDays}</dd>
              </div>
              <div>
                <dt>Reason:</dt>
                <dd>{sitExtensionReasons[sitExtension.requestReason]}</dd>
              </div>
              <div>
                <dt>Contractor remarks:</dt>
                <dd>{sitExtension.contractorRemarks}</dd>
              </div>
            </div>
            <Formik
              validationSchema={reviewSITExtensionSchema}
              onSubmit={(e) => onSubmit(sitExtension.id, e)}
              initialValues={{
                acceptExtension: 'yes',
                daysApproved: sitExtension.requestedDays.toString(),
                officeRemarks: '',
              }}
            >
              {({ isValid, values, setValues }) => {
                const handleNoSelection = (e) => {
                  if (e.target.value === 'no') {
                    setValues({
                      ...values,
                      daysApproved: '',
                      acceptExtension: 'no',
                    });
                  }
                };
                return (
                  <Form>
                    <FormGroup>
                      <Label>Accept request for extension?</Label>
                      <div>
                        <Field
                          as={Radio}
                          label="Yes"
                          id="acceptExtension"
                          name="acceptExtension"
                          value="yes"
                          title="Yes, accept extension"
                          type="radio"
                        />
                        <Field
                          as={Radio}
                          label="No"
                          id="denyExtension"
                          name="acceptExtension"
                          value="no"
                          title="No, deny extension"
                          type="radio"
                          onChange={handleNoSelection}
                        />
                      </div>
                    </FormGroup>
                    {values.acceptExtension === 'yes' && (
                      <MaskedTextField
                        name="daysApproved"
                        id="daysApproved"
                        label="Days approved"
                        mask="num"
                        blocks={{
                          num: {
                            mask: Number,
                            signed: false,
                            scale: 0,
                            thousandsSeparator: ',',
                          },
                        }}
                        lazy={false}
                        className={classnames(styles.ApprovedDaysInput, 'usa-input')}
                      />
                    )}
                    <Label htmlFor="officeRemarks">Office remarks</Label>
                    <Field
                      as={Textarea}
                      data-testid="officeRemarks"
                      label="No"
                      name="officeRemarks"
                      id="officeRemarks"
                    />
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

ReviewSITExtensionsModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  sitExtension: SITExtensionShape.isRequired,
};
export default ReviewSITExtensionsModal;

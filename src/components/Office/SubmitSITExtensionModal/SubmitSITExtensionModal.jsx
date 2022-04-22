import React from 'react';
import classnames from 'classnames';
import { Formik, Field } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button, Label, Textarea } from '@trussworks/react-uswds';

import styles from './SubmitSITExtensionModal.module.scss';

import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
import { DropdownInput } from 'components/form/fields';
import { sitExtensionReasons } from 'constants/sitExtensions';
import { dropdownInputOptions } from 'utils/formatters';

const reviewSITExtensionSchema = Yup.object().shape({
  requestReason: Yup.string().required('Required'),
  daysApproved: Yup.number()
    .min(1, 'Additional days approved must be greater than or equal to 1.')
    .required('Required'),
  officeRemarks: Yup.string().nullable(),
});

const SubmitSITExtensionModal = ({ onClose, onSubmit, summarySITComponent }) => {
  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.SubmitSITExtensionModal}>
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
                    <div className={styles.reasonDropdown}>
                      <DropdownInput
                        label="Reason for edit"
                        name="requestReason"
                        options={dropdownInputOptions(sitExtensionReasons)}
                      />
                    </div>
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

SubmitSITExtensionModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  summarySITComponent: PropTypes.node.isRequired,
};
export default SubmitSITExtensionModal;

import React from 'react';
import classnames from 'classnames';
import { Formik, Field, useField } from 'formik';
import PropTypes from 'prop-types';
import { Button, Label, Textarea } from '@trussworks/react-uswds';
import moment from 'moment';
import * as Yup from 'yup';

import styles from './EditSitEntryDateModal.module.scss';

import DataTableWrapper from 'components/DataTableWrapper/index';
import { DatePickerInput } from 'components/form/fields';
import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
import { formatDateForDatePicker, swaggerDateFormat } from 'shared/dates';
import RequiredAsterisk, { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

// datepickers that show the SIT entry dates - the previous one will be disabled
const SitEntryDateForm = ({ onChange }) => (
  <DatePickerInput
    name="sitEntryDate"
    label="New SIT entry date "
    id="sitEntryDate"
    onChange={onChange}
    showRequiredAsterisk
    required
  />
);
const DisabledSitEntryDateForm = () => (
  <DatePickerInput disabled name="prevSitEntryDate" label="Original SIT entry date" id="prevSitEntryDate" />
);

/**
 * @description This component contains the calendar pop outs and also sets the value when
 * the user changes the date in the proposed new SIT entry date change.
 */
const SitDatePickers = () => {
  // setting up a helper in order to update the values of form when the date is changed
  const entryDateHelper = useField({ name: 'sitEntryDate', id: 'sitEntryDate' })[2];
  const handleSitStartDateChange = (startDate) => {
    // Update form values
    entryDateHelper.setValue(startDate);
  };

  return (
    <div className={styles.formContainer}>
      {requiredAsteriskMessage}
      <div className={styles.formField}>
        <DisabledSitEntryDateForm />
      </div>
      <div className={styles.formField}>
        <SitEntryDateForm
          id="new-sit-entry-date"
          name="newSITEntryDate"
          required
          aria-required="true"
          onChange={(value) => {
            handleSitStartDateChange(value);
          }}
        />
      </div>
    </div>
  );
};

/**
 * @description This component contains a form that can be viewed from the MTO page
 * when a user clicks "Edit" next service items that are either
 * 1st day origin SIT || 1st day destination SIT
 */
const EditSitEntryDateModal = ({ onClose, onSubmit, serviceItem }) => {
  // setting initial values that requires some formatting for display requirements
  const initialValues = {
    sitEntryDate: formatDateForDatePicker(moment(serviceItem.sitEntryDate, swaggerDateFormat)),
    prevSitEntryDate: formatDateForDatePicker(moment(serviceItem.sitEntryDate, swaggerDateFormat)),
  };
  // right now the office remarks are just for show
  // TODO add change of SIT entry date to audit logs? Could be an enhancement
  // TODO I'm going to leave this here just in case
  const editSitEntryDateModalSchema = Yup.object().shape({
    officeRemarks: Yup.string().required('Required'),
    sitEntryDate: Yup.date()
      .transform((sitEntryDate) => {
        // in order to make sure the new value isn't the old value, we have to format it and check
        // this is because it is a Date value, but the prev value is formatted
        const formattedDate = formatDateForDatePicker(moment(sitEntryDate, swaggerDateFormat));
        if (formattedDate === initialValues.prevSitEntryDate) {
          throw new Yup.ValidationError(
            'New SIT entry date cannot be the same as the previous SIT entry date.',
            sitEntryDate,
            'sitEntryDate',
          );
        }

        // Return the formatted date if validation passes, or return an error message
        return sitEntryDate;
      })
      .notOneOf([Yup.ref('prevSitEntryDate')], 'SIT entry date cannot be the same as the previous entry date')
      .required('Required'),
  });

  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.UpdateSitEntryDateModal}>
          <ModalClose handleClick={() => onClose()} />
          <ModalTitle>
            <h2>Edit SIT Entry Date</h2>
          </ModalTitle>
          <div>
            <Formik
              onSubmit={(values) => onSubmit(serviceItem.id, values.sitEntryDate)}
              initialValues={{ ...initialValues }}
              validationSchema={editSitEntryDateModalSchema}
              initialTouched={{
                sitEntryDate: false,
                officeRemarks: true,
              }}
            >
              {({ isValid, setTouched, touched }) => {
                return (
                  <Form>
                    <DataTableWrapper
                      className={classnames('maxw-tablet', styles.sitDisplayForm)}
                      testID="sitExtensions"
                    >
                      <SitDatePickers serviceItem={serviceItem} />
                    </DataTableWrapper>
                    <Label htmlFor="officeRemarks" required>
                      <span required>
                        Office remarks <RequiredAsterisk />
                      </span>
                    </Label>
                    <Field
                      as={Textarea}
                      data-testid="officeRemarks"
                      label="No"
                      name="officeRemarks"
                      id="officeRemarks"
                      required
                    />
                    <ModalActions>
                      <Button
                        type="submit"
                        disabled={
                          (!isValid && (touched.sitEntryDate || touched.officeRemarks)) ||
                          !touched.sitEntryDate ||
                          !touched.officeRemarks
                        }
                        onClick={() => setTouched({ sitEntryDate: true, officeRemarks: true })}
                      >
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

EditSitEntryDateModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};
export default EditSitEntryDateModal;

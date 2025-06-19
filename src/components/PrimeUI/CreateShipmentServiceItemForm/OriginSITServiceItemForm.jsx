import * as Yup from 'yup';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';
import React from 'react';
import PropTypes from 'prop-types';

import formStyles from 'styles/form.module.scss';
import { requiredAddressSchema, ZIP_CODE_REGEX } from 'utils/validation';
import { formatDateForSwagger } from 'shared/dates';
import { formatAddressForPrimeAPI } from 'utils/formatters';
import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { DatePickerInput } from 'components/form/fields';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { ShipmentShape } from 'types/shipment';

const originSITValidationSchema = Yup.object().shape({
  reason: Yup.string().required('Required'),
  sitPostalCode: Yup.string().matches(ZIP_CODE_REGEX, 'Must be valid zip code').required('Required'),
  sitEntryDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .required('Required'),
  sitDepartureDate: Yup.date().typeError('Enter a complete date in DD MMM YYYY format (day, month, year).'),
  sitHHGActualOrigin: requiredAddressSchema,
});

const OriginSITServiceItemForm = ({ shipment, submission, handleCancel }) => {
  const initialValues = {
    moveTaskOrderID: shipment.moveTaskOrderID,
    mtoShipmentID: shipment.id,
    modelType: 'MTOServiceItemOriginSIT',
    reServiceCode: 'DOFSIT',
    reason: '',
    sitPostalCode: '',
    sitEntryDate: '',
    sitDepartureDate: '', // The Prime API is currently ignoring origin SIT departure date on creation
    sitHHGActualOrigin: {
      streetAddress1: '',
      streetAddress2: '',
      streetAddress3: '',
      city: '',
      state: '',
      postalCode: '',
      county: '',
    },
  };

  const onSubmit = (values) => {
    const { sitEntryDate, sitDepartureDate, sitHHGActualOrigin, ...serviceItemValues } = values;
    const body = {
      sitEntryDate: formatDateForSwagger(sitEntryDate),
      sitDepartureDate: sitDepartureDate ? formatDateForSwagger(sitDepartureDate) : null,
      sitHHGActualOrigin: sitHHGActualOrigin.streetAddress1 ? formatAddressForPrimeAPI(sitHHGActualOrigin) : null,
      ...serviceItemValues,
    };
    submission({ body });
  };

  return (
    <Formik initialValues={initialValues} validationSchema={originSITValidationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, ...formikProps }) => {
        return (
          <Form data-testid="originSITServiceItemForm" className={formStyles.form}>
            <input type="hidden" name="moveTaskOrderID" />
            <input type="hidden" name="mtoShipmentID" />
            <input type="hidden" name="modelType" />
            <input type="hidden" name="reServiceCode" />
            <TextField name="reason" id="reason" label="Reason" showRequiredAsterisk required />
            <MaskedTextField
              id="sitPostalCode"
              name="sitPostalCode"
              label="SIT postal code"
              mask="00000[{-}0000]"
              placeholder="62225"
              showRequiredAsterisk
              required
            />
            <DatePickerInput label="SIT entry Date" name="sitEntryDate" showRequiredAsterisk required />
            <DatePickerInput label="SIT departure Date" name="sitDepartureDate" />
            <h3 className={formStyles.sectionHeader}>SIT HHG actual origin address</h3>
            <AddressFields name="sitHHGActualOrigin" formikProps={formikProps} />
            <div className={formStyles.formActions}>
              <Button type="button" secondary onClick={handleCancel}>
                Cancel
              </Button>
              <Button onClick={handleSubmit} disabled={isSubmitting || !isValid} type="submit">
                Create service item
              </Button>
            </div>
          </Form>
        );
      }}
    </Formik>
  );
};

OriginSITServiceItemForm.propTypes = {
  shipment: ShipmentShape.isRequired,
  submission: PropTypes.func.isRequired,
};

export default OriginSITServiceItemForm;

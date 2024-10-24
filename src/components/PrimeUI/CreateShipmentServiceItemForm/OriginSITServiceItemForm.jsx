import * as Yup from 'yup';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';
import React from 'react';
import PropTypes from 'prop-types';

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

const OriginSITServiceItemForm = ({ shipment, submission }) => {
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
      city: '',
      state: '',
      postalCode: '',
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
      <Form data-testid="originSITServiceItemForm">
        <input type="hidden" name="moveTaskOrderID" />
        <input type="hidden" name="mtoShipmentID" />
        <input type="hidden" name="modelType" />
        <input type="hidden" name="reServiceCode" />
        <TextField name="reason" id="reason" label="Reason" />
        <MaskedTextField
          id="sitPostalCode"
          name="sitPostalCode"
          label="SIT postal code"
          mask="00000[{-}0000]"
          placeholder="62225"
        />
        <DatePickerInput label="SIT entry Date" name="sitEntryDate" />
        <DatePickerInput label="SIT departure Date" name="sitDepartureDate" />
        <AddressFields legend="SIT HHG actual origin" name="sitHHGActualOrigin" />
        <Button type="submit">Create service item</Button>
      </Form>
    </Formik>
  );
};

OriginSITServiceItemForm.propTypes = {
  shipment: ShipmentShape.isRequired,
  submission: PropTypes.func.isRequired,
};

export default OriginSITServiceItemForm;

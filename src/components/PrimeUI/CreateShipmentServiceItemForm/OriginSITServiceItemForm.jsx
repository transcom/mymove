import * as Yup from 'yup';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';
import React from 'react';
import PropTypes from 'prop-types';

import { requiredAddressSchema, ZIP_CODE_REGEX } from '../../../utils/validation';
import { formatDateForSwagger } from '../../../shared/dates';
import { formatAddressForPrimeAPI } from '../../../utils/formatters';
import { Form } from '../../form/Form';
import TextField from '../../form/fields/TextField';
import MaskedTextField from '../../form/fields/MaskedTextField';
import { DatePickerInput } from '../../form/fields/DatePickerInput';
import { AddressFields } from '../../form/AddressFields/AddressFields';
import { ShipmentShape } from '../../../types/shipment';

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
    sitDepartureDate: '',
    sitHHGActualOrigin: {
      street_address_1: '',
      street_address_2: '',
      city: '',
      state: '',
      postal_code: '',
    },
  };

  const onSubmit = (values) => {
    const { sitEntryDate, sitDepartureDate, sitHHGActualOrigin, ...serviceItemValues } = values;
    const body = {
      sitEntryDate: formatDateForSwagger(sitEntryDate),
      sitDepartureDate: sitDepartureDate ? formatDateForSwagger(sitDepartureDate) : null,
      sitHHGActualOrigin: sitHHGActualOrigin.street_address_1 ? formatAddressForPrimeAPI(sitHHGActualOrigin) : null,
      ...serviceItemValues,
    };
    submission({ body });
  };

  return (
    <Formik initialValues={initialValues} validationSchema={originSITValidationSchema} onSubmit={onSubmit}>
      <Form>
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

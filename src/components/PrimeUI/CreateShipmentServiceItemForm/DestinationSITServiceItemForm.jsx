import * as Yup from 'yup';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';
import React from 'react';
import PropTypes from 'prop-types';

import { addressSchema } from '../../../utils/validation';
import { formatDateForSwagger } from '../../../shared/dates';
import { formatAddressForPrimeAPI } from '../../../utils/formatters';
import { Form } from '../../form/Form';
import { DatePickerInput } from '../../form/fields/DatePickerInput';
import MaskedTextField from '../../form/fields/MaskedTextField';
import { AddressFields } from '../../form/AddressFields/AddressFields';
import { ShipmentShape } from '../../../types/shipment';

const destinationSITValidationSchema = Yup.object().shape({
  firstAvailableDeliveryDate1: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .required('Required'),
  timeMilitary1: Yup.string()
    .matches(/^(\d{4}Z)$/, 'Must be a valid military time (e.g. 1400Z)')
    .required('Required'),
  firstAvailableDeliveryDate2: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .required('Required'),
  timeMilitary2: Yup.string()
    .matches(/^(\d{4}Z)$/, 'Must be a valid military time (e.g. 1400Z)')
    .required('Required'),
  sitEntryDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .required('Required'),
  sitDepartureDate: Yup.date().typeError('Enter a complete date in DD MMM YYYY format (day, month, year).'),
  sitDestinationFinalAddress: addressSchema,
});

const DestinationSITServiceItemForm = ({ shipment, submission }) => {
  const initialValues = {
    moveTaskOrderID: shipment.moveTaskOrderID,
    mtoShipmentID: shipment.id,
    modelType: 'MTOServiceItemDestSIT',
    reServiceCode: 'DDFSIT',
    firstAvailableDeliveryDate1: '',
    timeMilitary1: '',
    firstAvailableDeliveryDate2: '',
    timeMilitary2: '',
    sitEntryDate: '',
    sitDepartureDate: '',
    sitDestinationFinalAddress: { street_address_1: '', street_address_2: '', city: '', state: '', postal_code: '' },
  };

  const onSubmit = (values) => {
    const {
      firstAvailableDeliveryDate1,
      firstAvailableDeliveryDate2,
      sitEntryDate,
      sitDepartureDate,
      sitDestinationFinalAddress,
      ...serviceItemValues
    } = values;
    const body = {
      firstAvailableDeliveryDate1: formatDateForSwagger(firstAvailableDeliveryDate1),
      firstAvailableDeliveryDate2: formatDateForSwagger(firstAvailableDeliveryDate2),
      sitEntryDate: formatDateForSwagger(sitEntryDate),
      sitDepartureDate: sitDepartureDate ? formatDateForSwagger(sitDepartureDate) : null,
      sitDestinationFinalAddress: sitDestinationFinalAddress.street_address_1
        ? formatAddressForPrimeAPI(sitDestinationFinalAddress)
        : null,
      ...serviceItemValues,
    };
    submission({ body });
  };

  return (
    <Formik initialValues={initialValues} validationSchema={destinationSITValidationSchema} onSubmit={onSubmit}>
      <Form>
        <input type="hidden" name="moveTaskOrderID" />
        <input type="hidden" name="mtoShipmentID" />
        <input type="hidden" name="modelType" />
        <input type="hidden" name="reServiceCode" />
        <DatePickerInput label="First available delivery date" name="firstAvailableDeliveryDate1" />
        <MaskedTextField
          id="timeMilitary1"
          name="timeMilitary1"
          label="First available delivery time"
          mask="0000{Z}"
          placeholder="1400Z"
        />
        <DatePickerInput label="Second available delivery date" name="firstAvailableDeliveryDate2" />
        <MaskedTextField
          id="timeMilitary1"
          name="timeMilitary2"
          label="Second available delivery time"
          mask="0000{Z}"
          placeholder="1400Z"
        />
        <DatePickerInput label="SIT entry date" name="sitEntryDate" />
        <DatePickerInput label="SIT departure date" name="sitDepartureDate" />
        <AddressFields legend="SIT destination final address" name="sitHHGActualOrigin" />
        <Button type="submit">Create service item</Button>
      </Form>
    </Formik>
  );
};

DestinationSITServiceItemForm.propTypes = {
  shipment: ShipmentShape.isRequired,
  submission: PropTypes.func.isRequired,
};

export default DestinationSITServiceItemForm;

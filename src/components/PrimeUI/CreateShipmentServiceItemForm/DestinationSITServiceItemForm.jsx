import * as Yup from 'yup';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';
import React from 'react';
import PropTypes from 'prop-types';

import { Form } from 'components/form/Form';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { formatDateForSwagger } from 'shared/dates';
import { formatAddressForPrimeAPI } from 'utils/formatters';
import { DatePickerInput } from 'components/form/fields';
import { ShipmentShape } from 'types/shipment';
import TextField from 'components/form/fields/TextField/TextField';
import Hint from 'components/Hint';

const destinationSITValidationSchema = Yup.object().shape({
  reason: Yup.string().required('Required'),
  firstAvailableDeliveryDate1: Yup.date().typeError('Enter a complete date in DD MMM YYYY format (day, month, year).'),
  timeMilitary1: Yup.string().matches(/^(\d{4}Z)$/, 'Must be a valid military time (e.g. 1400Z)'),
  firstAvailableDeliveryDate2: Yup.date().typeError('Enter a complete date in DD MMM YYYY format (day, month, year).'),
  timeMilitary2: Yup.string().matches(/^(\d{4}Z)$/, 'Must be a valid military time (e.g. 1400Z)'),
  sitEntryDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .required('Required'),
  sitDepartureDate: Yup.date().typeError('Enter a complete date in DD MMM YYYY format (day, month, year).'),
});

const DestinationSITServiceItemForm = ({ shipment, submission }) => {
  const initialValues = {
    moveTaskOrderID: shipment.moveTaskOrderID,
    mtoShipmentID: shipment.id,
    modelType: 'MTOServiceItemDestSIT',
    reServiceCode: 'DDFSIT',
    reason: '',
    firstAvailableDeliveryDate1: '',
    dateOfContact1: '',
    timeMilitary1: '',
    firstAvailableDeliveryDate2: '',
    dateOfContact2: '',
    timeMilitary2: '',
    sitEntryDate: '',
    sitDepartureDate: '',
    sitDestinationFinalAddress: { streetAddress1: '', streetAddress2: '', city: '', state: '', postalCode: '' },
  };

  const onSubmit = (values) => {
    const {
      firstAvailableDeliveryDate1,
      firstAvailableDeliveryDate2,
      sitEntryDate,
      sitDepartureDate,
      sitDestinationFinalAddress,
      timeMilitary1,
      timeMilitary2,
      dateOfContact1,
      dateOfContact2,
      ...serviceItemValues
    } = values;
    const body = {
      firstAvailableDeliveryDate1: formatDateForSwagger(firstAvailableDeliveryDate1),
      firstAvailableDeliveryDate2: formatDateForSwagger(firstAvailableDeliveryDate2),
      dateOfContact1: formatDateForSwagger(dateOfContact1),
      dateOfContact2: formatDateForSwagger(dateOfContact2),
      sitEntryDate: formatDateForSwagger(sitEntryDate),
      sitDepartureDate: sitDepartureDate ? formatDateForSwagger(sitDepartureDate) : null,
      sitDestinationFinalAddress: sitDestinationFinalAddress.streetAddress1
        ? formatAddressForPrimeAPI(sitDestinationFinalAddress)
        : null,
      timeMilitary1: timeMilitary1 || null,
      timeMilitary2: timeMilitary2 || null,
      ...serviceItemValues,
    };
    submission({ body });
  };

  return (
    <Formik initialValues={initialValues} validationSchema={destinationSITValidationSchema} onSubmit={onSubmit}>
      <Form data-testid="destinationSITServiceItemForm">
        <input type="hidden" name="moveTaskOrderID" />
        <input type="hidden" name="mtoShipmentID" />
        <input type="hidden" name="modelType" />
        <input type="hidden" name="reServiceCode" />
        <TextField label="Reason" name="reason" id="reason" />
        <DatePickerInput
          label="First available delivery date"
          name="firstAvailableDeliveryDate1"
          id="firstAvailableDeliveryDate1"
        />
        <DatePickerInput label="First date of attempted contact" name="dateOfContact1" id="dateOfContact1" />
        <MaskedTextField
          id="timeMilitary1"
          name="timeMilitary1"
          label="First time of attempted contact"
          mask="0000{Z}"
          placeholder="1400Z"
        />
        <DatePickerInput
          label="Second available delivery date"
          name="firstAvailableDeliveryDate2"
          id="firstAvailableDeliveryDate2"
        />
        <DatePickerInput label="Second date of attempted contact" name="dateOfContact2" id="dateOfContact2" />
        <MaskedTextField
          id="timeMilitary2"
          name="timeMilitary2"
          label="Second time of attempted contact"
          mask="0000{Z}"
          placeholder="1400Z"
        />
        <DatePickerInput label="SIT entry date" name="sitEntryDate" id="sitEntryDate" />
        <DatePickerInput label="SIT departure date" name="sitDepartureDate" id="sitDepartureDate" />
        <Hint data-testid="destinationSitInfo">
          The following service items will be created for domestic SIT: <br />
          DDFSIT (Domestic Destination 1st day SIT) <br />
          DDASIT (Domestic Destination additional days SIT) <br />
          DDDSIT (Domestic Destination SIT delivery) <br />
          DDSFSC (Domestic Destination SIT fuel surcharge) <br />
          <br />
          The following service items will be created for international SIT: <br />
          IDFSIT (International Destination 1st day SIT) <br />
          IDASIT (International Destination additional days SIT) <br />
          IDDSIT (International Destination SIT delivery) <br />
          IDSFSC (International Destination SIT fuel surcharge) <br />
          <br />
          <strong>NOTE:</strong> The above service items will use the current destination address of the shipment as
          their final destination address. Ensure the shipment address is accurate before creating these service items.
        </Hint>
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

import * as Yup from 'yup';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';
import React from 'react';
import PropTypes from 'prop-types';

// import { formatDateForSwagger } from '../../../shared/dates';
// import { formatAddressForPrimeAPI } from '../../../utils/formatters';

import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField';
import { ShipmentShape } from 'types/shipment';

const shuttleSITValidationSchema = Yup.object().shape({
  /*
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
  */
  reason: '',
});

const ShuttleSITServiceItemForm = ({ shipment, submission }) => {
  const initialValues = {
    moveTaskOrderID: shipment.moveTaskOrderID,
    mtoShipmentID: shipment.id,
    modelType: 'MTOServiceItemShuttle',
    // TODO: This Model type has two different `reServiceCode` associated with
    // it. This needs a Dropdown element associated with this form beyond what
    // was borrowed from the DestinationSITServiceItemForm that inspired this
    // file.
    reServiceCode: ['DOSHUT', 'DDSHUT'],
  };

  const onSubmit = (values) => {
    const {
      firstAvailableDeliveryDate1,
      /*
      firstAvailableDeliveryDate2,
      sitEntryDate,
      sitDepartureDate,
      sitDestinationFinalAddress,
      */
      ...serviceItemValues
    } = values;
    /*
    firstAvailableDeliveryDate1: formatDateForSwagger(firstAvailableDeliveryDate1),
    firstAvailableDeliveryDate2: formatDateForSwagger(firstAvailableDeliveryDate2),
    sitEntryDate: formatDateForSwagger(sitEntryDate),
    sitDepartureDate: sitDepartureDate ? formatDateForSwagger(sitDepartureDate) : null,
    sitDestinationFinalAddress: sitDestinationFinalAddress.streetAddress1
    ? formatAddressForPrimeAPI(sitDestinationFinalAddress)
    : null,
    ...serviceItemValues,
    */
    const body = {
      reServiceCode: 'DOSHUT',
      reason: '',
      estimatedWeight: 0,
      actualWeight: 0,
      ...serviceItemValues,
    };
    submission({ body });
  };

  return (
    <Formik initialValues={initialValues} validationSchema={shuttleSITValidationSchema} onSubmit={onSubmit}>
      <Form>
        <input type="hidden" name="moveTaskOrderID" />
        <input type="hidden" name="mtoShipmentID" />
        <input type="hidden" name="modelType" />
        <input type="hidden" name="reServiceCode" />
        <TextField name="reason" id="reason" label="Reason" />
        <Button type="submit">Create service item</Button>
      </Form>
    </Formik>
  );
};

ShuttleSITServiceItemForm.propTypes = {
  shipment: ShipmentShape.isRequired,
  submission: PropTypes.func.isRequired,
};

export default ShuttleSITServiceItemForm;

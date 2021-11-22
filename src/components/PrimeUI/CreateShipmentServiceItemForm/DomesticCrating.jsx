import * as Yup from 'yup';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';
import React from 'react';
import PropTypes from 'prop-types';

// import { formatDateForSwagger } from '../../../shared/dates';
// import { formatAddressForPrimeAPI } from '../../../utils/formatters';

import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField';
import { DropdownInput } from 'components/form/fields/DropdownInput';
import { ShipmentShape } from 'types/shipment';

const domesticShippingValidationSchema = Yup.object().shape({
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

const DomesticShippingServiceItemForm = ({ shipment, submission }) => {
  const initialValues = {
    moveTaskOrderID: shipment.moveTaskOrderID,
    mtoShipmentID: shipment.id,
    modelType: 'MTOServiceItemDomesticCrating',
    reServiceCode: ['DCRT', 'DUCRT'],
    itemLength: 0,
    itemWidth: 0,
    itemHeight: 0,
    crateLength: 0,
    crateWidth: 0,
    crateHeight: 0,
    reason: '',
    description: '',
  };

  const onSubmit = (values) => {
    const {
      itemLength,
      itemWidth,
      itemHeight,
      crateLength,
      crateWidth,
      crateHeight,
      reason,
      description,
      ...serviceItemValues
    } = values;

    const body = {
      reServiceCode: 'DCRT',
      item: {
        id: '',
        length: itemLength,
        width: itemWidth,
        height: itemHeight,
      },
      crate: {
        id: '',
        length: crateLength,
        width: crateWidth,
        height: crateHeight,
      },
      reason: '',
      description: '',
      ...serviceItemValues,
    };
    submission({ body });
  };

  return (
    <Formik initialValues={initialValues} validationSchema={domesticShippingValidationSchema} onSubmit={onSubmit}>
      <Form>
        <input type="hidden" name="moveTaskOrderID" />
        <input type="hidden" name="mtoShipmentID" />
        <input type="hidden" name="modelType" />
        <input type="hidden" name="reServiceCode" />
        <DropdownInput
          name="reServiceCode"
          id="reServiceCode"
          label="Service Code"
          options={[
            { value: 'DCRT', key: 'DCRT' },
            { value: 'DUCRT', key: 'DUCRT' },
          ]}
        />
        <TextField name="itemLength" id="itemLength" label="Item length" />
        <TextField name="itemWidth" id="itemWidth" label="Item width" />
        <TextField name="itemHeight" id="itemHeight" label="Item height" />
        <TextField name="crateLength" id="crateLength" label="Crate length" />
        <TextField name="crateWidth" id="crateWidth" label="Crate width" />
        <TextField name="crateHeight" id="crateHeight" label="Crate height" />
        <TextField name="description" id="description" label="Description" />
        <TextField name="reason" id="reason" label="Reason" />
        <Button type="submit">Create service item</Button>
      </Form>
    </Formik>
  );
};

DomesticShippingServiceItemForm.propTypes = {
  shipment: ShipmentShape.isRequired,
  submission: PropTypes.func.isRequired,
};

export default DomesticShippingServiceItemForm;

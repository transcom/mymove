import React, { useState } from 'react';
import { Formik } from 'formik';
import { Button, Dropdown, Label } from '@trussworks/react-uswds';

import { ShipmentShape } from '../../../types/shipment';
import TextField from '../../form/fields/TextField';
import { DatePickerInput } from '../../form/fields';
import Shipment from '../Shipment/Shipment';
import MaskedTextField from '../../form/fields/MaskedTextField';
import { Form } from '../../form';
import { AddressFields } from '../../form/AddressFields/AddressFields';

const serviceItemTypeOptions = (
  <>
    <option value="MTOServiceItemOriginSIT">Origin SIT</option>
    <option value="MTOServiceItemDestSIT">Destination SIT</option>
  </>
);

const originSITForm = (shipment) => {
  const initialValues = {
    moveTaskOrderID: shipment.moveTaskOrderID,
    mtoShipmentID: shipment.id,
    modelType: 'MTOServiceItemOriginSIT',
    reServiceCode: 'DOFSIT',
    reason: '',
    sitPostalCode: '',
    sitEntryDate: '',
    sitDepartureDate: '',
    sitHHGActualOrigin: {},
  };
  return (
    <Formik initialValues={initialValues}>
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
        <DatePickerInput label="SIT Entry Date" name="sitEntryDate" />
        <DatePickerInput label="SIT Departure Date" name="sitDepartureDate" />
        <AddressFields legend="SIT HHG actual origin" name="sitHHGActualOrigin" />
        <Button type="submit">Create service item</Button>
      </Form>
    </Formik>
  );
};

const destinationSITForm = (shipment) => {
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
    sitDestinationFinalAddress: {},
  };
  return (
    <Formik initialValues={initialValues}>
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

const populateServiceItemForm = (serviceItemType, shipment) => {
  switch (serviceItemType) {
    case 'MTOServiceItemOriginSIT':
      return originSITForm(shipment);
    case 'MTOServiceItemDestSIT':
      return destinationSITForm(shipment);
    default:
      return <></>;
  }
};

const CreateShipmentServiceItemForm = ({ shipment }) => {
  const [selectedServiceItemType, setSelectedServiceItemType] = useState('MTOServiceItemOriginSIT');

  const handleServiceItemTypeChange = (event) => {
    setSelectedServiceItemType(event.target.value);
  };

  return (
    <>
      <Shipment shipment={shipment} />
      <Label htmlFor="serviceItemType">Service item type</Label>
      <Dropdown id="serviceItemType" name="serviceItemType" onChange={handleServiceItemTypeChange}>
        {serviceItemTypeOptions}
      </Dropdown>
      {populateServiceItemForm(selectedServiceItemType, shipment)}
    </>
  );
};

CreateShipmentServiceItemForm.propTypes = {
  shipment: ShipmentShape.isRequired,
};

export default CreateShipmentServiceItemForm;

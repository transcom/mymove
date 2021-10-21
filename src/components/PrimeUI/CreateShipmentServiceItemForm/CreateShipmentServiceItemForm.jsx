import React, { useState } from 'react';
import { Formik } from 'formik';
import { Dropdown } from '@trussworks/react-uswds';

import { ShipmentShape } from '../../../types/shipment';
import TextField from '../../form/fields/TextField';
import { DatePickerInput } from '../../form/fields';
import Shipment from '../Shipment/Shipment';

const serviceItemTypeOptions = (
  <>
    <option value="MTOServiceItemOriginSIT">Origin SIT</option>
    <option value="MTOServiceItemDestSIT">Destination SIT</option>
  </>
);

const originSITForm = (
  <>
    <TextField name="reason" id="reason" label="Reason" />
  </>
);

const destinationSITForm = (
  <>
    <DatePickerInput label="First available delivery date" name="firstAvailableDeliveryDate1" />
  </>
);

const serviceItemFormLookup = {
  MTOServiceItemOriginSIT: originSITForm,
  MTOServiceItemDestSIT: destinationSITForm,
};

const CreateShipmentServiceItemForm = ({ shipment }) => {
  const [selectedServiceItemType, setSelectedServiceItemType] = useState('MTOServiceItemOriginSIT');

  const handleServiceItemTypeChange = (event) => {
    setSelectedServiceItemType(event.target.value);
  };

  return (
    <>
      <Shipment shipment={shipment} />
      <Dropdown name="serviceItemType" label="Service item type" onChange={handleServiceItemTypeChange}>
        {serviceItemTypeOptions}
      </Dropdown>
      <Formik>{serviceItemFormLookup[selectedServiceItemType]}</Formik>
    </>
  );
};

CreateShipmentServiceItemForm.propTypes = {
  shipment: ShipmentShape.isRequired,
};

export default CreateShipmentServiceItemForm;

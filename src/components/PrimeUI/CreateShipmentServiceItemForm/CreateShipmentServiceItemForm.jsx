import React, { useState } from 'react';
import { Dropdown, Label } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';

import { ShipmentShape } from '../../../types/shipment';
import { createServiceItemModelTypes } from '../../../constants/prime';

import styles from './CreateShipmentServiceItemForm.module.scss';
import OriginSITServiceItemForm from './OriginSITServiceItemForm';
import DestinationSITServiceItemForm from './DestinationSITServiceItemForm';

import Shipment from 'components/PrimeUI/Shipment/Shipment';

const CreateShipmentServiceItemForm = ({ shipment, createServiceItemMutation }) => {
  const { MTOServiceItemOriginSIT, MTOServiceItemDestSIT } = createServiceItemModelTypes;
  const [selectedServiceItemType, setSelectedServiceItemType] = useState(MTOServiceItemOriginSIT);

  const handleServiceItemTypeChange = (event) => {
    setSelectedServiceItemType(event.target.value);
  };

  return (
    <div className={styles.CreateShipmentServiceItemForm}>
      <Shipment shipment={shipment} />
      <Label htmlFor="serviceItemType">Service item type</Label>
      <Dropdown id="serviceItemType" name="serviceItemType" onChange={handleServiceItemTypeChange}>
        <>
          <option value={MTOServiceItemOriginSIT}>Origin SIT</option>
          <option value={MTOServiceItemDestSIT}>Destination SIT</option>
        </>
      </Dropdown>
      {selectedServiceItemType === MTOServiceItemOriginSIT && (
        <OriginSITServiceItemForm shipment={shipment} submission={createServiceItemMutation} />
      )}
      {selectedServiceItemType === MTOServiceItemDestSIT && (
        <DestinationSITServiceItemForm shipment={shipment} submission={createServiceItemMutation} />
      )}
    </div>
  );
};

CreateShipmentServiceItemForm.propTypes = {
  shipment: ShipmentShape.isRequired,
  createServiceItemMutation: PropTypes.func.isRequired,
};

export default CreateShipmentServiceItemForm;

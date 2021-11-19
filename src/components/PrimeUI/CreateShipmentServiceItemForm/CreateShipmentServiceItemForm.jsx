import React, { useState } from 'react';
import { Dropdown, Label } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';

import styles from './CreateShipmentServiceItemForm.module.scss';
import DestinationSITServiceItemForm from './DestinationSITServiceItemForm';
import OriginSITServiceItemForm from './OriginSITServiceItemForm';
import ShuttleSITServiceItemForm from './ShuttleSITServiceItemForm';

import { ShipmentShape } from 'types/shipment';
import { createServiceItemModelTypes } from 'constants/prime';
import Shipment from 'components/PrimeUI/Shipment/Shipment';

const CreateShipmentServiceItemForm = ({ shipment, createServiceItemMutation }) => {
  const { MTOServiceItemOriginSIT, MTOServiceItemDestSIT, MTOServiceItemShuttle, MTOServiceItemDomesticCrating } =
    createServiceItemModelTypes;
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
          <option value={MTOServiceItemShuttle}>Shuttle</option>
          <option value={MTOServiceItemDomesticCrating}>Domestic Crating</option>
        </>
      </Dropdown>
      {selectedServiceItemType === MTOServiceItemOriginSIT && (
        <OriginSITServiceItemForm shipment={shipment} submission={createServiceItemMutation} />
      )}
      {selectedServiceItemType === MTOServiceItemDestSIT && (
        <DestinationSITServiceItemForm shipment={shipment} submission={createServiceItemMutation} />
      )}
      {selectedServiceItemType === MTOServiceItemShuttle && (
        <ShuttleSITServiceItemForm shipment={shipment} submission={createServiceItemMutation} />
      )}
      {selectedServiceItemType === MTOServiceItemDomesticCrating && (
        /* <DestinationSITServiceItemForm shipment={shipment} submission={createServiceItemMutation} /> */
        <h1>Not Implemented Yet</h1>
      )}
    </div>
  );
};

CreateShipmentServiceItemForm.propTypes = {
  shipment: ShipmentShape.isRequired,
  createServiceItemMutation: PropTypes.func.isRequired,
};

export default CreateShipmentServiceItemForm;

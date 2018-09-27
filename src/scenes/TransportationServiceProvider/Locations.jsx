import React from 'react';
import { reduxForm, FormSection } from 'redux-form';

import { editablePanelify } from 'shared/EditablePanel';

import { addressElementDisplay, addressElementEdit } from 'shared/Address';

const LocationsDisplay = ({ shipment }) => {
  const {
    pickup_address,
    has_secondary_pickup_address,
    secondary_pickup_address,
  } = shipment;
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <span className="column-subhead">Pickup</span>
        {addressElementDisplay(pickup_address, 'Primary')}
        {has_secondary_pickup_address &&
          addressElementDisplay(secondary_pickup_address, 'Secondary')}
      </div>
    </React.Fragment>
  );
};

const LocationsEdit = props => {
  let { shipment } = props;
  let pickupProps = {
    swagger: props.addressSchema,
    values: props.pickupAddress,
  };
  let secondaryPickupProps = {
    swagger: props.addressSchema,
    values: props.secondaryPickupAddress,
  };
  let deliveryProps = {
    swagger: props.addressSchema,
    values: props.deliveryAddress,
  };
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <FormSection name="pickupAddress">
          <span className="column-subhead">Pickup</span>
          {addressElementEdit(pickupProps, 'Primary')}
          {shipment.has_secondary_pickup_address &&
            addressElementEdit(secondaryPickupProps, 'Secondary')}
        </FormSection>
      </div>
      <div className="editable-panel-column">
        <FormSection name="deliveryAddress">
          <span className="column-subhead">Delivery</span>
          {shipment.has_delivery_address &&
            addressElementEdit(deliveryProps, 'Primary')}
        </FormSection>
      </div>
    </React.Fragment>
  );
};

const formName = 'office_shipment_info_locations';
const editEnabled = false; // to remove the "Edit" button on panel header and disable editing

let LocationsPanel = editablePanelify(
  LocationsDisplay,
  LocationsEdit,
  editEnabled,
);
LocationsPanel = reduxForm({ form: formName })(LocationsPanel);

export default LocationsPanel;

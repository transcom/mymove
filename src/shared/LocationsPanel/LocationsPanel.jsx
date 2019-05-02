import { get } from 'lodash';
import React, { Fragment } from 'react';
import PropTypes from 'prop-types';
import { reduxForm } from 'redux-form';

import { editablePanelify } from 'shared/EditablePanel';
import { AddressElementDisplay, AddressElementEdit } from 'shared/Address';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { validateRequiredFields } from 'shared/JsonSchemaForm';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const LocationsDisplay = props => {
  // if they do not have a delivery address, default to the station's address info
  const deliveryAddress = props.shipment.has_delivery_address
    ? props.shipment.delivery_address
    : {
        city: props.newDutyStation.city,
        state: props.newDutyStation.state,
        postal_code: props.newDutyStation.postal_code,
      };
  const pickupAddress = props.shipment.pickup_address;
  const hasSecondaryPickupAddress = props.shipment.has_secondary_pickup_address;
  const secondaryPickupAddress = props.shipment.secondary_pickup_address;
  return (
    <Fragment>
      <div className="editable-panel-column">
        <span className="column-subhead">Pickup</span>
        {pickupAddress && <AddressElementDisplay address={pickupAddress} title="Primary" />}
        {hasSecondaryPickupAddress && <AddressElementDisplay address={secondaryPickupAddress} title="Additional" />}
      </div>
      <div className="editable-panel-column">
        <span className="column-subhead">Delivery</span>
        <AddressElementDisplay address={deliveryAddress} title="Primary" />
      </div>
    </Fragment>
  );
};

const LocationsEdit = props => {
  const { addressSchema, schema } = props;
  const newDutyStation = {
    city: props.newDutyStation.city,
    state: props.newDutyStation.state,
    postal_code: props.newDutyStation.postal_code,
  };
  const hasDeliveryAddress = get(props, 'formValues.has_delivery_address');
  const hasSecondaryPickupAddress = get(props, 'formValues.has_secondary_pickup_address');
  return (
    <Fragment>
      <div className="editable-panel-column">
        <span className="column-subhead">Pickup</span>
        <AddressElementEdit
          fieldName="pickup_address"
          schema={addressSchema}
          title="Primary address"
          zipPattern="USA"
          required
        />
        <SwaggerField
          className="radio-title"
          fieldName="has_secondary_pickup_address"
          swagger={schema}
          component={YesNoBoolean}
          title="Are there household goods at any other pickup location?"
        />
        {hasSecondaryPickupAddress && (
          <AddressElementEdit
            fieldName="secondary_pickup_address"
            schema={addressSchema}
            title="Additional address"
            zipPattern="USA"
          />
        )}
      </div>
      <div className="editable-panel-column">
        <span className="column-subhead">Delivery</span>
        <SwaggerField
          className="radio-title"
          fieldName="has_delivery_address"
          swagger={schema}
          component={YesNoBoolean}
          title="Do you know the delivery address at destination yet?"
        />
        {hasDeliveryAddress ? (
          <AddressElementEdit
            fieldName="delivery_address"
            schema={addressSchema}
            title="Primary address"
            zipPattern="USA"
          />
        ) : (
          <AddressElementDisplay address={newDutyStation} title="Delivery Primary (Duty Station)" />
        )}
      </div>
    </Fragment>
  );
};

const { shape, string, bool, object } = PropTypes;

const propTypes = {
  shipment: shape({
    pickup_address: shape({
      city: string.isRequired,
      postal_code: string.isRequired,
      state: string.isRequired,
      street_address_1: string.isRequired,
      street_address_2: string,
    }),
    has_secondary_pickup_address: bool.isRequired,
    secondary_pickup_address: shape({
      city: string.isRequired,
      postal_code: string.isRequired,
      state: string.isRequired,
      street_address_1: string.isRequired,
      street_address_2: string,
    }),
    has_delivery_address: bool.isRequired,
    delivery_address: shape({
      city: string.isRequired,
      postal_code: string.isRequired,
      state: string.isRequired,
      street_address_1: string,
      street_address_2: string,
    }),
  }),
  addressSchema: object.isRequired,
  schema: object.isRequired,
};

LocationsDisplay.propTypes = propTypes;
LocationsEdit.propTypes = propTypes;

const formName = 'shipment_locations';

let LocationsPanel = editablePanelify(LocationsDisplay, LocationsEdit);
LocationsPanel = reduxForm({
  form: formName,
  validate: validateRequiredFields,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(LocationsPanel);

LocationsPanel.propTypes = {
  initialValues: object.isRequired,
};

export { LocationsDisplay, LocationsPanel as default };

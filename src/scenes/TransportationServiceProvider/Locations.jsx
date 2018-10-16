import React, { Fragment } from 'react';
import PropTypes from 'prop-types';

import { editablePanelify } from 'shared/EditablePanel';
import { reduxForm, FormSection } from 'redux-form';

import { AddressElementDisplay, AddressElementEdit } from 'shared/Address';
import { validateRequiredFields } from 'shared/JsonSchemaForm';

const LocationsDisplay = ({
  deliveryAddress,
  shipment: {
    pickup_address: pickupAddress,
    has_secondary_pickup_address: hasSecondaryPickupAddress,
    secondary_pickup_address: secondaryPickupAddress,
  },
}) => (
  <Fragment>
    <div className="editable-panel-column">
      <span className="column-subhead">Pickup</span>
      <AddressElementDisplay address={pickupAddress} title="Primary" />
      {hasSecondaryPickupAddress && <AddressElementDisplay address={secondaryPickupAddress} title="Secondary" />}
    </div>
    <div className="editable-panel-column">
      <span className="column-subhead">Delivery</span>
      <AddressElementDisplay address={deliveryAddress} title="Primary" />
    </div>
  </Fragment>
);

const LocationsEdit = ({
  deliveryAddress,
  addressSchema,
  shipment: {
    pickup_address: pickupAddress,
    has_secondary_pickup_address: hasSecondaryPickupAddress,
    secondary_pickup_address: secondaryPickupAddress,
  },
}) => (
  <Fragment>
    <div className="editable-panel-column">
      <FormSection name="pickupAddress">
        <AddressElementEdit
          addressProps={{
            swagger: addressSchema,
            values: pickupAddress,
          }}
          title="Pickup Primary"
        />
      </FormSection>
      <FormSection name="secondaryPickupAddress">
        <AddressElementEdit
          addressProps={{
            swagger: addressSchema,
            values: secondaryPickupAddress,
          }}
          title="Pickup Secondary"
        />
      </FormSection>
    </div>
    <div className="editable-panel-column">
      <FormSection name="deliveryAddress">
        <AddressElementEdit
          addressProps={{
            swagger: addressSchema,
            values: deliveryAddress,
          }}
          title="Delivery Primary"
        />
      </FormSection>
    </div>
  </Fragment>
);

const { shape, string, bool, object } = PropTypes;

const propTypes = {
  deliveryAddress: shape({
    city: string.isRequired,
    postal_code: string.isRequired,
    state: string.isRequired,
    street_address_1: string,
    street_address_2: string,
    street_address_3: string,
  }).isRequired,
  shipment: shape({
    pickup_address: shape({
      city: string.isRequired,
      postal_code: string.isRequired,
      state: string.isRequired,
      street_address_1: string.isRequired,
      street_address_2: string,
      street_address_3: string,
    }),
    has_secondary_pickup_address: bool,
    secondary_pickup_address: shape({
      city: string.isRequired,
      postal_code: string.isRequired,
      state: string.isRequired,
      street_address_1: string.isRequired,
      street_address_2: string,
      street_address_3: string,
    }),
  }).isRequired,
  addressSchema: object.isRequired,
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

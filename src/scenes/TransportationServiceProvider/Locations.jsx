import { get } from 'lodash';
import React, { Fragment } from 'react';
import PropTypes from 'prop-types';

import { editablePanelify } from 'shared/EditablePanel';
import { reduxForm, FormSection } from 'redux-form';

import { AddressElementDisplay, AddressElementEdit } from 'shared/Address';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { validateRequiredFields } from 'shared/JsonSchemaForm';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const LocationsDisplay = props => {
  const { deliveryAddress } = props;
  const pickupAddress = props.shipment.pickup_address;
  const hasSecondaryPickupAddress = props.shipment.has_secondary_pickup_address;
  const secondaryPickupAddress = props.shipment.secondary_pickup_address;
  return (
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
};

const LocationsEdit = props => {
  const { deliveryAddress, addressSchema, schema } = props;
  const hasDeliveryAddress = get(props, 'formValues.hasDeliveryAddress', false);
  const hasSecondaryPickupAddress = get(props, 'formValues.hasSecondaryPickupAddress', false);
  const pickupAddress = props.shipment.pickup_address;
  const secondaryPickupAddress = props.shipment.secondary_pickup_address;
  return (
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
          <SwaggerField
            className="radio-title"
            fieldName="has_secondary_pickup_address"
            swagger={schema}
            component={YesNoBoolean}
          />
          {hasSecondaryPickupAddress && (
            <AddressElementEdit
              addressProps={{
                swagger: addressSchema,
                values: secondaryPickupAddress,
              }}
              title="Pickup Secondary"
            />
          )}
        </FormSection>
      </div>
      <div className="editable-panel-column">
        <FormSection name="deliveryAddress">
          <SwaggerField
            className="radio-title"
            fieldName="has_delivery_address"
            swagger={schema}
            component={YesNoBoolean}
          />
          {hasDeliveryAddress ? (
            <AddressElementEdit
              addressProps={{
                swagger: addressSchema,
                values: deliveryAddress,
              }}
              title="Delivery Primary"
            />
          ) : (
            <AddressElementDisplay address={deliveryAddress} title="Delivery Primary (Duty Station)" />
          )}
        </FormSection>
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
      street_address_3: string,
    }),
    has_secondary_pickup_address: bool.isRequired,
    secondary_pickup_address: shape({
      city: string.isRequired,
      postal_code: string.isRequired,
      state: string.isRequired,
      street_address_1: string.isRequired,
      street_address_2: string,
      street_address_3: string,
    }),
    has_delivery_address: bool.isRequired,
    delivery_address: shape({
      city: string.isRequired,
      postal_code: string.isRequired,
      state: string.isRequired,
      street_address_1: string,
      street_address_2: string,
      street_address_3: string,
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

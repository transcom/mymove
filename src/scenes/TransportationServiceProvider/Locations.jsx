import React, { Fragment } from 'react';
import PropTypes from 'prop-types';

import { editablePanelify } from 'shared/EditablePanel';

import { AddressElementDisplay } from 'shared/Address';

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
      {hasSecondaryPickupAddress && (
        <AddressElementDisplay
          address={secondaryPickupAddress}
          title="Secondary"
        />
      )}
    </div>
    <div className="editable-panel-column">
      <span className="column-subhead">Delivery</span>
      <AddressElementDisplay address={deliveryAddress} title="Primary" />
    </div>
  </Fragment>
);

const { shape, string, bool } = PropTypes;

LocationsDisplay.propTypes = {
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
};

const LocationsPanel = editablePanelify(LocationsDisplay, null, false);

export { LocationsDisplay, LocationsPanel as default };

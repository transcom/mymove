import React from 'react';
import PropTypes from 'prop-types';

import { editablePanelify } from 'shared/EditablePanel';

import { AddressElementDisplay } from 'shared/Address';

const LocationsDisplay = ({ shipment }) => {
  const {
    delivery_address,
    pickup_address,
    has_delivery_address,
    has_secondary_pickup_address,
    secondary_pickup_address,
    service_member,
  } = shipment;
  const { city, state, postal_code } = service_member.current_station.address;
  // if they do not have a delivery address, default to the station's address info
  const address = has_delivery_address
    ? delivery_address
    : { city, state, postal_code };
  return (
    <div className="editable-panel-column">
      <span className="column-subhead">Pickup</span>
      <AddressElementDisplay address={pickup_address} title="Primary" />
      {has_secondary_pickup_address && (
        <AddressElementDisplay
          address={secondary_pickup_address}
          title="Secondary"
        />
      )}
      <span className="column-subhead">Delivery</span>
      <AddressElementDisplay address={address} title="Primary" />
    </div>
  );
};

const { shape, string, number, bool } = PropTypes;

const address = shape({
  city: string.isRequired,
  country: string.isRequired,
  postal_code: string.isRequired,
  state: string.isRequired,
  street_address_1: string.isRequired,
  street_address_2: string,
  street_address_3: string,
});

LocationsDisplay.propTypes = {
  actual_delivery_date: string,
  actual_pickup_date: string,
  book_date: string,
  created_at: string,
  delivery_address: address,
  destination_gbloc: string,
  estimated_pack_days: number,
  estimated_transit_days: number,
  has_delivery_address: bool,
  has_secondary_pickup_address: bool,
  id: string,
  market: string,
  move: shape({
    cancel_reason: string,
    locator: string,
    orders_id: string.isRequired,
    selected_move_type: string,
    status: string.isRequired,
  }),
  pickup_address: address,
  progear_weight_estimate: number,
  requested_pickup_date: string,
  secondary_pickup_address: address,
  service_member: shape({
    current_station: shape({
      address,
    }),
    edipi: string,
    email_is_preferred: bool,
    first_name: string.isRequired,
    last_name: string.isRequired,
    personal_email: string.isRequired,
    telephone: string.isRequired,
  }),
  source_gbloc: string,
  spouse_progear_weight_estimate: number,
  status: string,
  traffic_distribution_list: shape({
    code_of_service: string,
    destination_region: string,
    source_rate_area: string,
  }),
  weight_estimate: number,
};

const LocationsPanel = editablePanelify(LocationsDisplay, null, false);

export { LocationsDisplay, LocationsPanel as default };

import React from 'react';
import Panel from 'shared/Panel';
import { AddressElementDisplay } from 'shared/Address';
import PropTypes from 'prop-types';

const LocationsPanel = ({ shipment }) => {
  const {
    pickup_address,
    has_secondary_pickup_address,
    secondary_pickup_address,
  } = shipment;
  return (
    <Panel>
      {({ isEditing }) => (
        <Panel.HalfRow>
          <Panel.Title editLabel="Edit">Locations</Panel.Title>
          {isEditing ? (
            <Panel.Content>
              <Panel.Header>Locations</Panel.Header>
              <Panel.Subheader>Primary Address</Panel.Subheader>
              <div>placeholder</div>
              <Panel.CancelButton>Cancel</Panel.CancelButton>
              <Panel.SaveButton>Save</Panel.SaveButton>
            </Panel.Content>
          ) : (
            <Panel.Content>
              <Panel.Subheader>Pickup</Panel.Subheader>
              <AddressElementDisplay address={pickup_address} title="Primary" />
              {has_secondary_pickup_address && (
                <AddressElementDisplay
                  address={secondary_pickup_address}
                  title="Secondary"
                />
              )}
            </Panel.Content>
          )}
        </Panel.HalfRow>
      )}
    </Panel>
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

LocationsPanel.propTypes = {
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

export default LocationsPanel;

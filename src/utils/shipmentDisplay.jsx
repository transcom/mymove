/*  camelcase */
import React from 'react';

import { shipmentOptionLabels } from 'shared/constants';

export function formatAddress(address) {
  const { street_address_1, city, state, postal_code } = address;
  return (
    <>
      {street_address_1 && <>{street_address_1},&nbsp;</>}
      {city ? `${city}, ${state} ${postal_code}` : postal_code}
    </>
  );
}

export function formatCustomerDestination(destinationLocation, destinationZIP) {
  return destinationLocation ? (
    <>
      {destinationLocation.street_address_1} {destinationLocation.street_address_2}
      <br />
      {destinationLocation.city}, {destinationLocation.state} {destinationLocation.postal_code}
    </>
  ) : (
    destinationZIP
  );
}

export const getShipmentTypeLabel = (shipmentType) => shipmentOptionLabels.find((l) => l.key === shipmentType)?.label;

export function formatPaymentRequestAddressString(pickupAddress, destinationAddress) {
  if (pickupAddress && destinationAddress) {
    return `${pickupAddress.city}, ${pickupAddress.state} ${pickupAddress.postal_code} to ${destinationAddress.city}, ${destinationAddress.state} ${destinationAddress.postal_code}`;
  }
  if (pickupAddress && !destinationAddress) {
    return `${pickupAddress.city}, ${pickupAddress.state} ${pickupAddress.postal_code} to TBD`;
  }
  if (!pickupAddress && destinationAddress) {
    return `TBD to ${destinationAddress.city}, ${destinationAddress.state} ${destinationAddress.postal_code}`;
  }
  return ``;
}

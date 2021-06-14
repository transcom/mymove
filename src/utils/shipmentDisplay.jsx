/* eslint-disable camelcase */
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React from 'react';

import { shipmentOptionLabels } from 'shared/constants';
import { shipmentStatuses, shipmentModificationTypes } from 'constants/shipments';

export function formatAddress(address) {
  const { street_address_1, city, state, postal_code } = address;
  return (
    <>
      {street_address_1 && <>{street_address_1},&nbsp;</>}
      {city ? `${city}, ${state} ${postal_code}` : postal_code}
    </>
  );
}

export function formatAgent(agent) {
  const { firstName, lastName, phone, email } = agent;
  return (
    <>
      <div>
        {firstName} {lastName}
      </div>
      {phone && <div>{phone}</div>}
      {email && <div>{email}</div>}
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
    return (
      <>
        {pickupAddress.city}, {pickupAddress.state} {pickupAddress.postal_code} <FontAwesomeIcon icon="arrow-right" />{' '}
        {destinationAddress.city}, {destinationAddress.state} {destinationAddress.postal_code}
      </>
    );
  }
  if (pickupAddress && !destinationAddress) {
    return (
      <>
        {pickupAddress.city}, {pickupAddress.state} {pickupAddress.postal_code} <FontAwesomeIcon icon="arrow-right" />{' '}
        TBD
      </>
    );
  }
  if (!pickupAddress && destinationAddress) {
    return (
      <>
        TBD <FontAwesomeIcon icon="arrow-right" /> {destinationAddress.city}, {destinationAddress.state}{' '}
        {destinationAddress.postal_code}
      </>
    );
  }
  return ``;
}

export function formatPaymentRequestReviewAddressString(address) {
  if (address) {
    return `${address.city}, ${address.state} ${address.postal_code}`;
  }
  return '';
}

export function getShipmentModificationType(shipment) {
  if (shipment.status === shipmentStatuses.CANCELED) {
    return shipmentModificationTypes.CANCELED;
  }

  if (shipment.diversion === true) {
    return shipmentModificationTypes.DIVERSION;
  }

  return '';
}

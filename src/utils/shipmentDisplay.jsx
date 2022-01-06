/* eslint-disable camelcase */
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React from 'react';

import { LOA_TYPE, shipmentOptionLabels } from 'shared/constants';
import { shipmentStatuses, shipmentModificationTypes } from 'constants/shipments';

export function formatAddress(address) {
  const { streetAddress1, city, state, postalCode } = address;
  return (
    <>
      {streetAddress1 && <>{streetAddress1},&nbsp;</>}
      {city ? `${city}, ${state} ${postalCode}` : postalCode}
    </>
  );
}

export function retrieveTAC(tacType, ordersLOA) {
  switch (tacType) {
    case LOA_TYPE.HHG:
      return ordersLOA.tac;
    case LOA_TYPE.NTS:
      return ordersLOA.ntsTac;
    default:
      return ordersLOA.tac;
  }
}

export function retrieveSAC(sacType, ordersLOA) {
  switch (sacType) {
    case LOA_TYPE.HHG:
      return ordersLOA.sac;
    case LOA_TYPE.NTS:
      return ordersLOA.ntsSac;
    default:
      return ordersLOA.sac;
  }
}

export function formatAccountingCode(accountingCode, accountingCodeType) {
  return String(accountingCode).concat(' (', accountingCodeType, ')');
}

// Display street address 1, street address 2, city, state, and zip
// for Prime API Prime Simulator UI shipment
export function formatPrimeAPIShipmentAddress(address) {
  return address?.id ? (
    <>
      {address.streetAddress1 && <>{address.streetAddress1},&nbsp;</>}
      {address.streetAddress2 && <>{address.streetAddress2},&nbsp;</>}
      {address.city ? `${address.city}, ${address.state} ${address.postalCode}` : address.postalCode}
    </>
  ) : (
    ''
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
      {destinationLocation.streetAddress1} {destinationLocation.streetAddress2}
      <br />
      {destinationLocation.city}, {destinationLocation.state} {destinationLocation.postalCode}
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
        {pickupAddress.city}, {pickupAddress.state} {pickupAddress.postalCode} <FontAwesomeIcon icon="arrow-right" />{' '}
        {destinationAddress.city}, {destinationAddress.state} {destinationAddress.postalCode}
      </>
    );
  }
  if (pickupAddress && !destinationAddress) {
    return (
      <>
        {pickupAddress.city}, {pickupAddress.state} {pickupAddress.postalCode} <FontAwesomeIcon icon="arrow-right" />{' '}
        TBD
      </>
    );
  }
  if (!pickupAddress && destinationAddress) {
    return (
      <>
        TBD <FontAwesomeIcon icon="arrow-right" /> {destinationAddress.city}, {destinationAddress.state}{' '}
        {destinationAddress.postalCode}
      </>
    );
  }
  return ``;
}

export function formatPaymentRequestReviewAddressString(address) {
  if (address) {
    return `${address.city}, ${address.state} ${address.postalCode}`;
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

  return undefined;
}

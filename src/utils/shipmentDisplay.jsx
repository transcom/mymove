/* eslint-disable camelcase */
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React from 'react';

import { LOA_TYPE, shipmentOptionLabels } from 'shared/constants';
import { shipmentStatuses, shipmentModificationTypes } from 'constants/shipments';
import affiliations from 'content/serviceMemberAgencies';

export function formatAddress(pickupAddress) {
  const { streetAddress1, streetAddress2, streetAddress3, city, state, postalCode } = pickupAddress;

  if (streetAddress1 === 'n/a') {
    return city ? `${city}, ${state} ${postalCode}` : postalCode;
  }
  return (
    <>
      {streetAddress1 && <>{streetAddress1},&nbsp;</>}
      {streetAddress2 && <>{streetAddress2},&nbsp;</>}
      {streetAddress3 && <>{streetAddress3},&nbsp;</>}
      {city ? `${city}, ${state} ${postalCode}` : postalCode}
    </>
  );
}

export function formatTwoLineAddress(address) {
  const { streetAddress1, streetAddress2, streetAddress3, city, state, postalCode } = address;

  return (
    <address data-testid="two-line-address">
      {streetAddress1 && `${streetAddress1},`}
      {streetAddress2 && ` ${streetAddress2}`}
      {streetAddress3 && ` ${streetAddress3}`}
      <br />
      {city ? `${city}, ${state} ${postalCode}` : postalCode}
    </address>
  );
}

/**
 * @description This function is used to format the address in the
 * EditSitAddressChangeForm component. It specifically uses the `<span>`
 * elements to be able to make each line of an address be set to `display:
 * block;` in the CSS to match the design.
 * @see ServiceItemUpdateModal / EditSITAddressChangeForm
 * */
export function formatAddressForSitAddressChangeForm(address) {
  const { streetAddress1, streetAddress2, streetAddress3, city, state, postalCode } = address;
  return (
    <address data-testid="SitAddressChangeDisplay">
      {streetAddress1 && <span data-testid="AddressLine">{streetAddress1},</span>}
      {streetAddress2 && <span data-testid="AddressLine">{streetAddress2},</span>}
      {streetAddress3 && <span data-testid="AddressLine">{streetAddress3},</span>}
      <span data-testid="AddressLine">{city ? `${city}, ${state} ${postalCode}` : postalCode}</span>
    </address>
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
      {address.streetAddress3 && <>{address.streetAddress3},&nbsp;</>}
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
      {destinationLocation.streetAddress1} {destinationLocation.streetAddress2} {destinationLocation.streetAddress3}
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

/**
 * @description This function is used to format the address in the
 * ShipmentAddresses and PaymentRequestReview components.
 * It displays only the city, state and postal code.
 * */
export function formatCityStateAndPostalCode(address) {
  if (address) {
    return `${address.city}, ${address.state} ${address.postalCode}`;
  }
  return '';
}

/**
 * @description This function is used to format the port in the
 * ShipmentAddresses component.
 * It displays only the port code, port name, city, state and zip code.
 * */
export function formatPortInfo(port) {
  if (port) {
    return `${port.portCode} - ${port.portName}\n${port.city}, ${port.state} ${port.zip}`;
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

/**
 * @description Returns whether the service member is affilated a branch that is allowed to choose their PPM closeout office.
 * @param {string} affiliation - String representing the service member's branch (i.e. "AIR_FORCE") (these string constants are defined in src/content/serviceMemberAgencies.js)
 * @returns {boolean} - True if member is affiliated with any of the branches allowed to choose their PPM closeout office.
 */
export function canChoosePPMLocation(affiliation) {
  return (
    affiliation === affiliations.AIR_FORCE ||
    affiliation === affiliations.ARMY ||
    affiliation === affiliations.SPACE_FORCE
  );
}

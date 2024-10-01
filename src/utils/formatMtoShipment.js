import { isEmpty } from 'lodash';
import moment from 'moment';

import { MTOAgentType, SHIPMENT_TYPES } from 'shared/constants';
import { parseDate } from 'shared/dates';
import { formatDelimitedNumber, parseSwaggerDate } from 'utils/formatters';
import { roleTypes } from 'constants/userRoles';
import { LOCATION_TYPES } from 'types/sitStatusShape';
import { boatShipmentTypes } from 'constants/shipments';

const formatDateForSwagger = (date) => {
  const parsedDate = parseDate(date);
  if (parsedDate) {
    return moment(parsedDate).format('YYYY-MM-DD');
  }
  return '';
};

function formatAgentForDisplay(agent) {
  const agentCopy = { ...agent };
  return agentCopy;
}

function formatAgentForAPI(agent) {
  const agentCopy = { ...agent };
  Object.keys(agentCopy).forEach((key) => {
    const sanitizedKey = `${key}`;
    if (
      // These fields are readOnly, so we don't want to send them in requests
      sanitizedKey === 'updatedAt' ||
      sanitizedKey === 'createdAt' ||
      sanitizedKey === 'mtoShipmentID'
    ) {
      delete agentCopy[sanitizedKey];
    }
  });
  return agentCopy;
}

export function formatStorageFacilityForAPI(storageFacility) {
  const storageFacilityCopy = { ...storageFacility };
  Object.keys(storageFacilityCopy).forEach((key) => {
    const sanitizedKey = `${key}`;
    if (storageFacilityCopy[sanitizedKey] === '') {
      delete storageFacilityCopy[sanitizedKey];
    } else if (
      // These fields are readOnly so we don't want to send them in requests
      sanitizedKey === 'eTag'
    ) {
      delete storageFacilityCopy[sanitizedKey];
    }
  });
  return storageFacilityCopy;
}

export function removeEtag(obj) {
  const { eTag, ...rest } = obj;
  return rest;
}

export function formatAddressForAPI(address) {
  const formattedAddress = address;

  if (formattedAddress.state) {
    formattedAddress.state = formattedAddress.state?.toUpperCase();
    delete formattedAddress.id;
    return formattedAddress;
  }

  return undefined;
}

const emptyAgentShape = {
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
};

const emptyAddressShape = {
  streetAddress1: '',
  streetAddress2: '',
  city: '',
  state: '',
  postalCode: '',
};

export function formatPpmShipmentForDisplay({ counselorRemarks = '', ppmShipment = {}, closeoutOffice = {} }) {
  const displayValues = {
    expectedDepartureDate: ppmShipment.expectedDepartureDate,
    pickup: {
      address: ppmShipment.pickupAddress || emptyAddressShape,
    },
    destination: {
      address: ppmShipment.destinationAddress || emptyAddressShape,
    },
    secondaryPickup: {
      address: { ...emptyAddressShape },
    },
    secondaryDestination: {
      address: { ...emptyAddressShape },
    },
    hasSecondaryPickup: ppmShipment.hasSecondaryPickupAddress ? 'true' : 'false',
    hasSecondaryDestination: ppmShipment.hasSecondaryDestinationAddress ? 'true' : 'false',

    tertiaryPickup: {
      address: { ...emptyAddressShape },
    },
    tertiaryDestination: {
      address: { ...emptyAddressShape },
    },
    hasTertiaryPickup: ppmShipment.hasTertiaryPickupAddress ? 'true' : 'false',
    hasTertiaryDestination: ppmShipment.hasTertiaryDestinationAddress ? 'true' : 'false',

    sitExpected: !!ppmShipment.sitExpected,
    sitLocation: ppmShipment.sitLocation ?? LOCATION_TYPES.DESTINATION,
    sitEstimatedWeight: (ppmShipment.sitEstimatedWeight || '').toString(),
    sitEstimatedEntryDate: ppmShipment.sitEstimatedEntryDate,
    sitEstimatedDepartureDate: ppmShipment.sitEstimatedDepartureDate,

    estimatedWeight: (ppmShipment.estimatedWeight || '').toString(),
    hasProGear: !!ppmShipment.hasProGear,
    proGearWeight: (ppmShipment.proGearWeight || '').toString(),
    spouseProGearWeight: (ppmShipment.spouseProGearWeight || '').toString(),

    estimatedIncentive: ppmShipment.estimatedIncentive,
    advanceRequested: ppmShipment.hasRequestedAdvance ?? false,
    advanceStatus: ppmShipment.advanceStatus,
    advance: (ppmShipment.advanceAmountRequested / 100 || '').toString(),
    closeoutOffice,
    counselorRemarks,
  };

  if (ppmShipment.hasSecondaryPickupAddress) {
    displayValues.secondaryPickup.address = { ...emptyAddressShape, ...ppmShipment.secondaryPickupAddress };
  }

  if (ppmShipment.hasSecondaryDestinationAddress) {
    displayValues.secondaryDestination.address = { ...emptyAddressShape, ...ppmShipment.secondaryDestinationAddress };
  }

  if (ppmShipment.hasTertiaryPickupAddress) {
    displayValues.tertiaryPickup.address = { ...emptyAddressShape, ...ppmShipment.tertiaryPickupAddress };
  }

  if (ppmShipment.hasTertiaryDestinationAddress) {
    displayValues.tertiaryDestination.address = { ...emptyAddressShape, ...ppmShipment.tertiaryDestinationAddress };
  }

  return displayValues;
}

/**
 * formatMtoShipmentForDisplay converts mtoShipment data from the format API calls expect to the template format
 * @param {*} mtoShipment - (see MtoShipmentShape)
 */
export function formatMtoShipmentForDisplay({
  agents,
  shipmentType,
  requestedPickupDate,
  pickupAddress,
  requestedDeliveryDate,
  destinationAddress,
  customerRemarks,
  counselorRemarks,
  moveTaskOrderID,
  secondaryPickupAddress,
  secondaryDeliveryAddress,
  tertiaryPickupAddress,
  tertiaryDeliveryAddress,

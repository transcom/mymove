import { isEmpty } from 'lodash';
import moment from 'moment';

import { MTOAgentType } from 'shared/constants';
import { parseDate } from 'shared/dates';
import { formatDelimitedNumber, parseSwaggerDate } from 'utils/formatters';
import { roleTypes } from 'constants/userRoles';
import { LOCATION_TYPES } from 'types/sitStatusShape';

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
    pickupPostalCode: ppmShipment.pickupPostalCode || '',
    secondPickupPostalCode: ppmShipment.secondaryPickupPostalCode || '',
    destinationPostalCode: ppmShipment.destinationPostalCode || '',
    secondDestinationPostalCode: ppmShipment.secondaryDestinationPostalCode || '',

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
  ntsRecordedWeight,
  tacType,
  sacType,
  serviceOrderNumber,
  storageFacility,
  usesExternalVendor,
  userRole,
  destinationType,
}) {
  const displayValues = {
    shipmentType,
    moveTaskOrderID,
    customerRemarks: customerRemarks || '',
    counselorRemarks: counselorRemarks || '',
    pickup: {
      requestedDate: '',
      address: { ...emptyAddressShape },
      agent: { ...emptyAgentShape },
    },
    delivery: {
      requestedDate: '',
      address: { ...emptyAddressShape },
      agent: { ...emptyAgentShape },
    },
    secondaryPickup: {
      address: { ...emptyAddressShape },
    },
    secondaryDelivery: {
      address: { ...emptyAddressShape },
    },
    hasDeliveryAddress: 'no',
    hasSecondaryPickup: 'no',
    hasSecondaryDelivery: 'no',
    ntsRecordedWeight,
    tacType,
    sacType,
    serviceOrderNumber,
    usesExternalVendor,
  };

  if (agents) {
    const receivingAgent = agents.find((agent) => agent.agentType === 'RECEIVING_AGENT');
    const releasingAgent = agents.find((agent) => agent.agentType === 'RELEASING_AGENT');

    if (receivingAgent) {
      const formattedAgent = formatAgentForDisplay(receivingAgent);
      if (Object.keys(formattedAgent).length) {
        displayValues.delivery.agent = { ...emptyAgentShape, ...formattedAgent };
      }
    }
    if (releasingAgent) {
      const formattedAgent = formatAgentForDisplay(releasingAgent);
      if (Object.keys(formattedAgent).length) {
        displayValues.pickup.agent = { ...emptyAgentShape, ...formattedAgent };
      }
    }
  }

  if (pickupAddress) {
    displayValues.pickup.address = { ...emptyAddressShape, ...pickupAddress };
  }

  if (requestedPickupDate) {
    displayValues.pickup.requestedDate = parseSwaggerDate(requestedPickupDate);
  }

  if (secondaryPickupAddress) {
    displayValues.secondaryPickup.address = { ...emptyAddressShape, ...secondaryPickupAddress };
    displayValues.hasSecondaryPickup = 'yes';
  }

  if (destinationAddress) {
    displayValues.delivery.address = { ...emptyAddressShape, ...destinationAddress };
    displayValues.hasDeliveryAddress = 'yes';
  }

  if (destinationType) {
    displayValues.destinationType = destinationType;
  }

  if (secondaryDeliveryAddress) {
    displayValues.secondaryDelivery.address = { ...emptyAddressShape, ...secondaryDeliveryAddress };
    displayValues.hasSecondaryDelivery = 'yes';
  }

  if (requestedDeliveryDate) {
    displayValues.delivery.requestedDate = parseSwaggerDate(requestedDeliveryDate);
  }

  if (storageFacility) {
    displayValues.storageFacility = {
      ...storageFacility,
      address: {
        ...emptyAddressShape,
        ...(storageFacility?.address || {}),
      },
    };
  }

  if (userRole === roleTypes.TOO && usesExternalVendor === undefined) {
    // Vendor defaults to the Prime
    displayValues.usesExternalVendor = false;
  }

  return displayValues;
}

export function formatPpmShipmentForAPI(formValues) {
  let ppmShipmentValues = {
    expectedDepartureDate: formatDateForSwagger(formValues.expectedDepartureDate),
    pickupPostalCode: formValues.pickupPostalCode,
    secondaryPickupPostalCode: formValues.secondPickupPostalCode || undefined,
    destinationPostalCode: formValues.destinationPostalCode,
    secondaryDestinationPostalCode: formValues.secondDestinationPostalCode || undefined,
    sitExpected: !!formValues.sitExpected,
    estimatedWeight: Number(formValues.estimatedWeight || '0'),
    hasProGear: !!formValues.hasProGear,
    hasRequestedAdvance: formValues.advanceRequested,
    advanceStatus: formValues.advanceStatus,
  };

  if (formValues.hasProGear) {
    ppmShipmentValues = {
      ...ppmShipmentValues,
      proGearWeight: formValues.proGearWeight ? Number(formValues.proGearWeight) : undefined,
      spouseProGearWeight: formValues.spouseProGearWeight ? Number(formValues.spouseProGearWeight) : undefined,
    };
  }

  if (formValues.sitExpected) {
    ppmShipmentValues = {
      ...ppmShipmentValues,
      sitLocation: formValues.sitLocation,
      sitEstimatedWeight: formValues.sitEstimatedWeight ? Number(formValues.sitEstimatedWeight || '0') : undefined,
      sitEstimatedEntryDate: formValues.sitEstimatedEntryDate
        ? formatDateForSwagger(formValues.sitEstimatedEntryDate)
        : undefined,
      sitEstimatedDepartureDate: formValues.sitEstimatedDepartureDate
        ? formatDateForSwagger(formValues.sitEstimatedDepartureDate)
        : undefined,
    };
  }

  if (formValues.advanceRequested) {
    ppmShipmentValues = {
      ...ppmShipmentValues,
      advanceAmountRequested: formValues.advance ? Number(formValues.advance || '0') * 100 : undefined,
    };
  }

  return {
    shipmentType: 'PPM',
    counselorRemarks: formValues.counselorRemarks === undefined ? undefined : formValues.counselorRemarks,
    ppmShipment: ppmShipmentValues,
  };
}

/**
 * formatMtoShipmentForAPI converts mtoShipment data from the template format to the format API calls expect
 * @param {*} param - unnamed object representing various mtoShipment data parts
 */
export function formatMtoShipmentForAPI({
  moveId,
  shipmentType,
  pickup,
  delivery,
  customerRemarks,
  counselorRemarks,
  hasSecondaryPickup,
  secondaryPickup,
  hasSecondaryDelivery,
  secondaryDelivery,
  ntsRecordedWeight,
  tacType,
  sacType,
  serviceOrderNumber,
  storageFacility,
  usesExternalVendor,
  destinationType,
}) {
  const formattedMtoShipment = {
    moveTaskOrderID: moveId,
    shipmentType,
    customerRemarks,
    counselorRemarks,
    agents: [],
    destinationType,
    tacType,
    sacType,
  };

  if (pickup?.requestedDate && pickup.requestedDate !== '') {
    formattedMtoShipment.requestedPickupDate = formatDateForSwagger(pickup.requestedDate);
    formattedMtoShipment.pickupAddress = formatAddressForAPI(pickup.address);

    if (pickup.agent) {
      const formattedAgent = formatAgentForAPI(pickup.agent);
      if (!isEmpty(formattedAgent)) {
        formattedMtoShipment.agents.push({ ...formattedAgent, agentType: MTOAgentType.RELEASING });
      }
    }
  }

  if (delivery?.requestedDate && delivery.requestedDate !== '') {
    formattedMtoShipment.requestedDeliveryDate = formatDateForSwagger(delivery.requestedDate);

    if (delivery.address) {
      formattedMtoShipment.destinationAddress = formatAddressForAPI(delivery.address);
    }

    if (destinationType) {
      formattedMtoShipment.destinationType = destinationType;
    }

    if (delivery.agent) {
      const formattedAgent = formatAgentForAPI(delivery.agent);
      if (!isEmpty(formattedAgent)) {
        formattedMtoShipment.agents.push({ ...formattedAgent, agentType: MTOAgentType.RECEIVING });
      }
    }
  }
  formattedMtoShipment.hasSecondaryPickupAddress = hasSecondaryPickup;

  if (hasSecondaryPickup && secondaryPickup?.address) {
    formattedMtoShipment.secondaryPickupAddress = formatAddressForAPI(secondaryPickup.address);
  }

  formattedMtoShipment.hasSecondaryDeliveryAddress = hasSecondaryDelivery;
  if (hasSecondaryDelivery && secondaryDelivery?.address) {
    formattedMtoShipment.secondaryDeliveryAddress = formatAddressForAPI(secondaryDelivery.address);
  }

  if (!formattedMtoShipment.agents?.length) {
    formattedMtoShipment.agents = undefined;
  }

  if (ntsRecordedWeight) {
    formattedMtoShipment.ntsRecordedWeight = formatDelimitedNumber(ntsRecordedWeight);
  }

  if (serviceOrderNumber) {
    formattedMtoShipment.serviceOrderNumber = serviceOrderNumber;
  }

  if (storageFacility?.address) {
    const sanitizedStorageFacility = formatStorageFacilityForAPI(storageFacility);
    formattedMtoShipment.storageFacility = {
      ...sanitizedStorageFacility,
      address: removeEtag(formatAddressForAPI(sanitizedStorageFacility.address)),
    };
  }

  if (usesExternalVendor !== undefined) {
    formattedMtoShipment.usesExternalVendor = usesExternalVendor;
  }

  return formattedMtoShipment;
}

export function getMtoShipmentLabel({ context }) {
  if (!context) return {};
  const mtoShipmentLabels = {};
  if (context[0].shipment_type) {
    mtoShipmentLabels.shipment_type = context[0].shipment_type;
  }
  if (context[0].shipment_id_abbr) {
    mtoShipmentLabels.shipment_id_display = context[0].shipment_id_abbr.toUpperCase();
  }
  if (context[0].name) {
    mtoShipmentLabels.service_item_name = context[0].name;
  }
  return mtoShipmentLabels;
}

export default {
  formatMtoShipmentForAPI,
  formatMtoShipmentForDisplay,
  formatAddressForAPI,
  formatStorageFacilityForAPI,
  removeEtag,
  getMtoShipmentLabel,
};

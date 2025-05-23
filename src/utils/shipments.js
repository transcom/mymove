// The PPM shipment creation is a multi-step flow so it's possible to get in a state with missing
// information and get to the review screen in an incomplete state from creating another shipment

import { PPM_TYPES, SHIPMENT_OPTIONS } from '../shared/constants';

import { expenseTypes } from 'constants/ppmExpenseTypes';

// on the move. hasRequestedAdvance is the last required field that would mean they're finished.
export function isPPMShipmentComplete(mtoShipment) {
  if (mtoShipment?.ppmShipment?.hasRequestedAdvance != null) {
    return true;
  }
  return false;
}

// isPPMAboutInfoComplete - checks if all the "About your ppm" fields have data in them.
export function isPPMAboutInfoComplete(ppmShipment) {
  const hasBaseRequiredFields = [
    'actualMoveDate',
    'pickupAddress',
    'destinationAddress',
    'w2Address',
    'hasReceivedAdvance',
  ].every((fieldName) => ppmShipment[fieldName] !== null);

  if (hasBaseRequiredFields) {
    if (
      !ppmShipment.hasReceivedAdvance ||
      (ppmShipment.advanceAmountReceived !== null && ppmShipment.advanceAmountReceived > 0)
    ) {
      return true;
    }
  }

  return false;
}

// isWeightTicketComplete - checks that the required fields for a weight ticket have valid data
// to check if the weight ticket can be considered complete. For the purposes of this function,
// any data is enough to consider some fields valid.
export function isWeightTicketComplete(weightTicket) {
  const hasValidEmptyWeight = weightTicket.emptyWeight != null && weightTicket.emptyWeight >= 0;

  const hasTrailerDocUpload = weightTicket.proofOfTrailerOwnershipDocument.uploads.length > 0;
  const needsTrailerUpload = weightTicket.ownsTrailer && weightTicket.trailerMeetsCriteria;
  const trailerNeedsMet = needsTrailerUpload ? hasTrailerDocUpload : true;

  return !!(
    weightTicket.vehicleDescription &&
    hasValidEmptyWeight &&
    weightTicket.emptyDocument.uploads.length > 0 &&
    weightTicket.fullWeight > 0 &&
    weightTicket.fullDocument.uploads.length > 0 &&
    trailerNeedsMet
  );
}

// hasCompletedAllWeightTickets - checks if every weight ticket has been completed.
// Returns false if there are no weight tickets, or if any of them are incomplete.
export function hasCompletedAllWeightTickets(weightTickets, ppmType) {
  // PPM-SPRs don't have weight tickets
  if (ppmType === PPM_TYPES.SMALL_PACKAGE) {
    return true;
  }
  if (!weightTickets?.length) {
    return false;
  }

  return !!weightTickets?.every(isWeightTicketComplete);
}

export default isPPMShipmentComplete;

// isExpenseComplete - checks that the required fields for an expense have valid data
// to check if the expense can be considered complete. For the purposes of this function,
// any data is enough to consider some fields valid.
export function isExpenseComplete(expense) {
  const hasADocumentUpload = expense.document.uploads.length > 0;
  const hasValidSITDates =
    expense.movingExpenseType !== expenseTypes.STORAGE || (expense.sitStartDate && expense.sitEndDate);
  const requiresDescription = expense.movingExpenseType !== expenseTypes.SMALL_PACKAGE;
  return !!(
    (requiresDescription ? expense.description : true) &&
    expense.movingExpenseType &&
    expense.amount &&
    hasADocumentUpload &&
    hasValidSITDates
  );
}

// hasCompletedAllExpenses - checks if expense ticket has been completed.
// Returns true if expenses are not defined or there are none, false if any of them are incomplete.
export function hasCompletedAllExpenses(expenses) {
  if (!expenses?.length) {
    return true;
  }

  return !!expenses?.every(isExpenseComplete);
}

export function isProGearComplete(proGear) {
  const hasADocumentUpload = proGear.document.uploads.length > 0;
  const hasAnOwner = proGear.belongsToSelf !== null;
  return !!(proGear.weight && proGear.description && hasADocumentUpload && hasAnOwner);
}

export function hasCompletedAllProGear(proGear) {
  if (!proGear?.length) {
    return true;
  }
  return !!proGear?.every(isProGearComplete);
}

export function isPPM(shipment) {
  return shipment?.shipmentType === SHIPMENT_OPTIONS.PPM;
}

export function isPPMOnly(mtoShipments) {
  if (!mtoShipments?.length) {
    return false;
  }
  return !!mtoShipments?.every(isPPM);
}
export function isBoatShipmentComplete(mtoShipment) {
  return mtoShipment?.requestedPickupDate;
}

export function isMobileHomeShipmentComplete(mtoShipment) {
  return mtoShipment?.requestedPickupDate;
}

export function hasIncompleteWeightTicket(weightTickets) {
  if (!weightTickets?.length) {
    return false;
  }

  return !weightTickets?.every(isWeightTicketComplete);
}

export const blankAddress = {
  address: {
    streetAddress1: '',
    streetAddress2: '',
    streetAddress3: '',
    city: '',
    state: '',
    postalCode: '',
    usPostRegionCitiesID: '',
  },
};

const updateAddressToggle = (setValues, fieldName, value, fieldKey, fieldValue) => {
  if (fieldName === 'hasSecondaryPickup' && value === 'false') {
    // HHG
    setValues((prevValues) => ({
      ...prevValues,
      [fieldName]: value,
      [fieldKey]: fieldValue,
      hasTertiaryPickup: 'false',
      tertiaryPickup: {
        address: blankAddress,
      },
    }));
  } else if (fieldName === 'hasDeliveryAddress' && value === 'false') {
    // HHG
    setValues((prevValues) => ({
      ...prevValues,
      [fieldName]: value,
      [fieldKey]: fieldValue,
      hasSecondaryDelivery: 'false',
      secondaryDelivery: {
        address: blankAddress,
      },
      hasTertiaryDelivery: 'false',
      tertiaryDelivery: {
        address: blankAddress,
      },
    }));
  } else if (fieldName === 'hasSecondaryDelivery' && value === 'false') {
    // HHG
    setValues((prevValues) => ({
      ...prevValues,
      [fieldName]: value,
      [fieldKey]: fieldValue,
      hasTertiaryDelivery: 'false',
      tertiaryDelivery: {
        address: blankAddress,
      },
    }));
  } else if (fieldName === 'hasSecondaryPickupAddress' && value === 'false') {
    // PPM
    setValues((prevValues) => ({
      ...prevValues,
      [fieldName]: value,
      [fieldKey]: fieldValue,
      hasTertiaryPickupAddress: 'false',
      tertiaryPickupAddress: {
        address: blankAddress,
      },
    }));
  } else if (fieldName === 'hasSecondaryDestinationAddress' && value === 'false') {
    // PPM
    setValues((prevValues) => ({
      ...prevValues,
      [fieldName]: value,
      [fieldKey]: fieldValue,
      hasTertiaryDestinationAddress: 'false',
      tertiaryDestinationAddress: {
        address: blankAddress,
      },
    }));
  } else {
    setValues((prevValues) => ({
      ...prevValues,
      [fieldName]: value,
      [fieldKey]: value === 'false' ? fieldValue : { ...prevValues[fieldKey] },
    }));
  }
};

export const handleAddressToggleChange = (e, values, setValues, newDutyLocationAddress) => {
  const { name, value } = e.target;

  const fieldMap = {
    hasSecondaryPickup: { key: 'secondaryPickup', updateValue: { blankAddress } },
    hasSecondaryPickupAddress: { key: 'secondaryPickupAddress', updateValue: { blankAddress } },
    hasTertiaryPickup: { key: 'tertiaryPickup', updateValue: { blankAddress } },
    hasTertiaryPickupAddress: { key: 'tertiaryPickupAddress', updateValue: { blankAddress } },
    hasDeliveryAddress: {
      key: 'delivery',
      updateValue: {
        ...values.delivery,
        address: {
          streetAddress1: 'N/A',
          city: newDutyLocationAddress.city,
          state: newDutyLocationAddress.state,
          postalCode: newDutyLocationAddress.postalCode,
          county: newDutyLocationAddress.county,
          usPostRegionCitiesID: newDutyLocationAddress.usPostRegionCitiesID,
        },
      },
    },
    hasSecondaryDelivery: { key: 'secondaryDelivery', updateValue: { blankAddress } },
    hasSecondaryDestination: { key: 'secondaryDestination', updateValue: { blankAddress } },
    hasSecondaryDestinationAddress: { key: 'secondaryDestinationAddress', updateValue: { blankAddress } },
    hasTertiaryDelivery: { key: 'tertiaryDelivery', updateValue: { blankAddress } },
    hasTertiaryDestination: { key: 'tertiaryDestination', updateValue: { blankAddress } },
    hasTertiaryDestinationAddress: { key: 'tertiaryDestinationAddress', updateValue: { blankAddress } },
  };

  if (fieldMap[name]) {
    updateAddressToggle(setValues, name, value, fieldMap[name].key, fieldMap[name].updateValue);
  }
};

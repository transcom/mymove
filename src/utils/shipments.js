// The PPM shipment creation is a multi-step flow so it's possible to get in a state with missing
// information and get to the review screen in an incomplete state from creating another shipment

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
    'actualPickupPostalCode',
    'actualDestinationPostalCode',
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
export function hasCompletedAllWeightTickets(weightTickets) {
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
  return !!(
    expense.description &&
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

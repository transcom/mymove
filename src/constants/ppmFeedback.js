import { expenseTypeLabels } from './ppmExpenseTypes';
import { PPM_DOCUMENT_STATUS } from './ppms';

import { formatCents, formatCustomerDate, formatWeight, formatYesNoInputValue, toDollarString } from 'utils/formatters';

export const FEEDBACK_DOCUMENT_TYPES = {
  WEIGHT: 'Trip',
  PRO_GEAR: 'Set',
  MOVING_EXPENSE: 'Receipt',
};

const feedbackDisplayHelperTrip = (documentSet) => {
  return documentSet?.some(
    (doc) =>
      (doc.status !== null && doc.status !== PPM_DOCUMENT_STATUS.APPROVED) ||
      (doc.submittedEmptyWeight != null && doc.submittedEmptyWeight !== doc.emptyWeight) ||
      (doc.submittedFullWeight != null && doc.submittedFullWeight !== doc.fullWeight) ||
      (doc.submittedOwnsTrailer != null && doc.submittedOwnsTrailer !== doc.ownsTrailer) ||
      (doc.submittedTrailerMeetsCriteria != null && doc.submittedTrailerMeetsCriteria !== doc.trailerMeetsCriteria),
  );
};

const feedbackDisplayHelperProGear = (documentSet) => {
  return documentSet?.some(
    (doc) =>
      (doc.status !== null && doc.status !== PPM_DOCUMENT_STATUS.APPROVED) ||
      (doc.submittedBelongsToSelf != null && doc.submittedBelongsToSelf !== doc.belongsToSelf) ||
      (doc.submittedHasWeightTickets != null && doc.submittedHasWeightTickets !== doc.hasWeightTickets) ||
      (doc.submittedWeight != null && doc.submittedWeight !== doc.weight),
  );
};

const feedbackDisplayHelperExpense = (documentSet) => {
  return documentSet?.some(
    (doc) =>
      (doc.status !== null && doc.status !== PPM_DOCUMENT_STATUS.APPROVED) ||
      (doc.submittedAmount != null && doc.submittedAmount !== doc.amount) ||
      (doc.submittedDescription != null && doc.submittedDescription !== doc.description) ||
      (doc.submittedMovingExpenseType != null && doc.submittedMovingExpenseType !== doc.movingExpenseType) ||
      (doc.submittedSitEndDate != null && doc.submittedSitEndDate !== doc.sitEndDate) ||
      (doc.submittedSitStartDate != null && doc.submittedSitStartDate !== doc.sitStartDate),
  );
};

// feedback should only NOT be visible if all ppm documents were accepted without edits
// each of the above helper functions returns true if any document is NOT approved
// or if any customer submitted value does NOT equal the final value
export const isFeedbackAvailable = (ppmShipment) => {
  if (!ppmShipment) return false;
  if (feedbackDisplayHelperTrip(ppmShipment?.weightTickets)) return true;
  if (feedbackDisplayHelperProGear(ppmShipment?.proGearWeightTickets)) return true;
  if (feedbackDisplayHelperExpense(ppmShipment?.movingExpenses)) return true;
  return false;
};

const getExpenseType = (label) => {
  return expenseTypeLabels[label];
};

// helper function to handle label with boolean value
const formatProGearLabel = (belongsToSelf) => {
  return belongsToSelf ? 'Pro-Gear' : 'Spouse Pro-Gear';
};

// templates for feedback items are stored as arrays to allow for ordering
// key - corresponds to the key in the document object
// label - the label for the value in the UI
// format - use if the value in the document needs formatting
// secondaryKey - used for submitted_ columns so we can track changes made by closeout SC
const FEEDBACK_TRIP_TEMPLATE = [
  { key: 'vehicleDescription', label: 'Vehicle description: ' },
  {
    key: 'emptyWeight',
    label: 'Empty: ',
    format: (weight) => formatWeight(weight),
    secondaryKey: 'submittedEmptyWeight',
  },
  {
    key: 'fullWeight',
    label: 'Full: ',
    format: (weight) => formatWeight(weight),
    secondaryKey: 'submittedFullWeight',
  },
  { key: 'tripWeight', label: 'Trip weight: ', format: (weight) => formatWeight(weight) },
  {
    key: 'ownsTrailer',
    label: 'Trailer: ',
    format: (bool) => formatYesNoInputValue(bool),
    secondaryKey: 'submittedOwnsTrailer',
  },
  {
    key: 'trailerMeetsCriteria',
    label: 'Trailer meets criteria: ',
    format: (bool) => formatYesNoInputValue(bool),
    secondaryKey: 'submittedtrailerMeetsCriteria',
  },
  { key: 'status' },
];

const FEEDBACK_SET_TEMPLATE = [
  {
    key: 'belongsToSelf',
    label: '',
    format: (bool) => formatProGearLabel(bool),
    secondaryKey: 'submittedBelongsToSelf',
  },
  { key: 'description', label: 'Description: ' },
  { key: 'weight', label: 'Weight: ', format: (weight) => formatWeight(weight), secondaryKey: 'submittedWeight' },
  {
    key: 'hasWeightTickets',
    label: 'Weight tickets: ',
    format: (bool) => formatYesNoInputValue(bool),
    secondaryKey: 'submittedhasWeightTickets',
  },
  { key: 'status' },
];

const FEEDBACK_RECEIPT_TEMPLATE = [
  {
    key: 'movingExpenseType',
    label: 'Type: ',
    format: (expenseType) => getExpenseType(expenseType),
    secondaryKey: 'submittedMovingExpenseType',
  },
  { key: 'description', label: 'Description: ', secondaryKey: 'submittedDescription' },
  {
    key: 'amount',
    label: 'Amount: ',
    format: (amount) => toDollarString(formatCents(amount)),
    secondaryKey: 'submittedAmount',
  },
  {
    key: 'sitStartDate',
    label: 'SIT start date: ',
    format: (date) => formatCustomerDate(date),
    secondaryKey: 'submittedSitStartDate',
  },
  {
    key: 'sitEndDate',
    label: 'SIT end date: ',
    format: (date) => formatCustomerDate(date),
    secondaryKey: 'submittedSitEndDate',
  },
  { key: 'status' },
];

export const FEEDBACK_TEMPLATES = {
  Trip: FEEDBACK_TRIP_TEMPLATE,
  Set: FEEDBACK_SET_TEMPLATE,
  Receipt: FEEDBACK_RECEIPT_TEMPLATE,
};

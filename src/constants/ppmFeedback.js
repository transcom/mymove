import { expenseTypeLabels } from './ppmExpenseTypes';
import ppms from './ppms';

import { formatCents, formatCustomerDate, formatWeight, formatYesNoInputValue, toDollarString } from 'utils/formatters';

const feedbackDisplayHelperTrip = (documentSet) => {
  return documentSet?.some(
    (doc) =>
      (doc.status !== ppms.APPROVED && doc.status !== null) ||
      doc.submittedEmptyWeight !== doc.emptyWeight ||
      doc.submittedFullWeight !== doc.fullWeight ||
      doc.submittedOwnsTrailer !== doc.ownsTrailer ||
      doc.submittedTrailerMeetsCriteria !== doc.trailerMeetsCriteria,
  );
};

const feedbackDisplayHelperProGear = (documentSet) => {
  return documentSet?.some(
    (doc) =>
      (doc.status !== ppms.APPROVED && doc.status !== null) ||
      doc.submittedBelongsToSelf !== doc.belongsToSelf ||
      doc.submittedHasWeightTickets !== doc.hasWeightTickets ||
      doc.submittedWeight !== doc.weight,
  );
};

const feedbackDisplayHelperExpense = (documentSet) => {
  return documentSet?.some(
    (doc) =>
      (doc.status !== ppms.APPROVED && doc.status !== null) ||
      doc.submittedAmount !== doc.amount ||
      doc.submittedDescription !== doc.description ||
      doc.submittedMovingExpenseType !== doc.movingExpenseType ||
      doc.submittedSitEndDate !== doc.sitEndDate ||
      doc.submittedSitStartDate !== doc.sitStartDate,
  );
};

// feedback should only be visible if all ppm documents were accepted without edits
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

const formatProGearLabel = (belongsToSelf) => {
  return belongsToSelf ? 'Pro-Gear' : 'Spouse Pro-Gear';
};

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

export const FEEDBACK_DOCUMENT_TYPES = {
  WEIGHT: 'Trip',
  PRO_GEAR: 'Set',
  MOVING_EXPENSE: 'Receipt',
};

// refactor?
export const getFeedbackTemplate = (type) => {
  let template;
  switch (type) {
    case 'Trip':
      template = FEEDBACK_TRIP_TEMPLATE;
      break;
    case 'Set':
      template = FEEDBACK_SET_TEMPLATE;
      break;
    case 'Receipt':
      template = FEEDBACK_RECEIPT_TEMPLATE;
      break;
    default:
      break;
  }
  return template;
};

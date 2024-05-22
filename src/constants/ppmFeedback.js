import { expenseTypeLabels } from './ppmExpenseTypes';

import { formatCents, formatCustomerDate, formatWeight, formatYesNoInputValue, toDollarString } from 'utils/formatters';

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
  { key: 'ownsTrailer', label: 'Trailer: ', format: (bool) => formatYesNoInputValue(bool) },
  { key: 'status' },
];

const FEEDBACK_SET_TEMPLATE = [
  { key: 'belongsToSelf', label: '', format: (bool) => formatProGearLabel(bool) },
  { key: 'description', label: 'Description: ' },
  { key: 'weight', label: 'Weight: ', format: (weight) => formatWeight(weight), secondaryKey: 'submittedWeight' },
  { key: 'status' },
];

const FEEDBACK_RECEIPT_TEMPLATE = [
  { key: 'movingExpenseType', label: 'Type: ', format: (expenseType) => getExpenseType(expenseType) },
  { key: 'description', label: 'Description: ' },
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

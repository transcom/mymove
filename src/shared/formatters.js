import { isFinite } from 'lodash';
import moment from 'moment';
import path from 'path';
import { ORDERS_TYPE_DETAILS_OPTIONS, ORDERS_TYPE_OPTIONS } from 'constants/orders';
import { PAYMENT_REQUEST_STATUS_LABELS } from 'constants/paymentRequestStatus';
import { SERVICE_COUNSELING_MOVE_STATUS_OPTIONS } from 'constants/queues';

export function formatNumber(num) {
  if (!isFinite(num)) {
    return '';
  }

  return num.toLocaleString();
}

// Format a number of cents into a string, e.g. $12,345.67
export function formatCents(cents) {
  if (!isFinite(cents)) {
    return '';
  }

  return (cents / 100).toLocaleString(undefined, {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });
}

// Format base quantity as cents
export function formatBaseQuantityAsDollars(baseQuantity) {
  return formatCents(baseQuantity / 100);
}

// Format a base quantity into a user-friendly number string, e.g. 167000 -> "16.7000"
export function formatFromBaseQuantity(baseQuantity) {
  if (!isFinite(baseQuantity)) {
    return '';
  }

  return (baseQuantity / 10000).toLocaleString(undefined, {
    minimumFractionDigits: 4,
    maximumFractionDigits: 4,
  });
}

// Format a base quantity into a user-friendly number, e.g. 167000 -> 16.7
export function convertFromBaseQuantity(baseQuantity) {
  if (!isFinite(baseQuantity)) {
    return null;
  }

  return baseQuantity / 10000;
}

// Format user-entered base quantity into base quantity, e.g. 16.7000 -> 167000
export function formatToBaseQuantity(baseQuantity) {
  baseQuantity = parseFloat(String(baseQuantity).replace(',', '')) * 10000;

  return baseQuantity;
}

// Format a thousandth of an inch into an inch, e.g. 16700 -> 16.7
export function convertFromThousandthInchToInch(thousandthInch) {
  if (!isFinite(thousandthInch)) {
    return null;
  }

  return thousandthInch / 1000;
}

// Format a dimensions object length, width and height to inches
export function formatToDimensionsInches(dimensions) {
  if (!dimensions) {
    return undefined;
  }

  dimensions.length = convertFromThousandthInchToInch(dimensions.length);
  dimensions.width = convertFromThousandthInchToInch(dimensions.width);
  dimensions.height = convertFromThousandthInchToInch(dimensions.height);

  return dimensions;
}

// Format dimensions object length, width and height to base dimensions
export function formatDimensionsToThousandthInches(dimensions) {
  if (!dimensions) {
    return;
  }

  dimensions.length = formatToThousandthInches(dimensions.length);
  dimensions.width = formatToThousandthInches(dimensions.width);
  dimensions.height = formatToThousandthInches(dimensions.height);
  return;
}

// Format user-entered dimension into base dimension, e.g. 15.25 -> 15250
export function formatToThousandthInches(val) {
  return parseFloat(String(val).replace(',', '')) * 1000;
}

export function formatCentsRange(min, max) {
  if (!isFinite(min) || !isFinite(max)) {
    return '';
  }

  return `$${formatCents(min)} - ${formatCents(max)}`;
}

// Format a date, include its time and timezone, e.g. 03-Jan-2018 21:23 ET
export function formatDateTimeWithTZ(date) {
  if (!date) return undefined;

  // This gets us a date string that includes the browser timezone
  // e.g. Mon Apr 22 2019 09:08:10 GMT-0500 (Central Daylight Time)
  // If this looks a bit strange, it's a workaround for IE11 not
  // supporting the timeZoneName: 'short' option in Date.toLocaleString
  const newDateString = String(moment(date).toDate());
  const longZone = newDateString.substring(newDateString.lastIndexOf('(') + 1, newDateString.lastIndexOf(')'));
  let shortZone = longZone
    .split(' ')
    .map((word) => word[0])
    .join('');

  // Converting timezones like CDT and EST to CT and ET
  if (shortZone.length === 3 && shortZone !== 'UTC') {
    shortZone = shortZone.slice(0, 1) + shortZone.slice(2, 3);
  }

  return moment(date, moment.ISO_8601, true).format('DD-MMM-YY HH:mm') + ` ${shortZone}`;
}

export function formatTimeAgo(date) {
  if (!date) return undefined;

  return moment(date)
    .fromNow()
    .replace('minute', 'min')
    .replace(/a min\s/, '1 min ');
}

// truncate a number and return appropiate decimal places... (watch out for negitive numbers: floor(-5.1) === -6)
// see test for examples of how this works
export const truncateNumber = (num, decimalPlaces = 0) => {
  if (!num) return num;

  const floatNum = parseFloat(num).toFixed(4);
  const scale = Math.pow(10, decimalPlaces);
  const truncatedNbr = Math.floor(floatNum * scale) / scale;
  return truncatedNbr.toFixed(decimalPlaces).toString();
};

// adds commas to numberString w/o removeing .0000 from the end of the string or rounding
export const addCommasToNumberString = (numOrString, decimalPlaces = 0) => {
  if (!numOrString || numOrString === '0') {
    numOrString = (0).toFixed(decimalPlaces);
  }

  const str = numOrString.toString();
  const [wholeNum, decimalNum] = str.split('.');
  const wholeNumInt = parseInt(wholeNum);
  if (decimalNum) {
    return `${wholeNumInt.toLocaleString()}.${decimalNum}`;
  }
  return wholeNumInt.toLocaleString();
};

export const serviceCounselingMoveStatusLabel = (status) => {
  return SERVICE_COUNSELING_MOVE_STATUS_OPTIONS.find((option) => option.value === `${status}`)?.label || status;
};

export const ordersTypeReadable = (ordersType, missingText) => {
  if (!ordersType) {
    return missingText;
  }
  return ORDERS_TYPE_OPTIONS[`${ordersType}`] || ordersType;
};

export const ordersTypeDetailReadable = (ordersTypeDetail, missingText) => {
  if (!ordersTypeDetail) {
    return missingText;
  }
  return ORDERS_TYPE_DETAILS_OPTIONS[`${ordersTypeDetail}`] || ordersTypeDetail;
};

export const paymentRequestStatusReadable = (paymentRequestStatus) => {
  return PAYMENT_REQUEST_STATUS_LABELS[`${paymentRequestStatus}`] || paymentRequestStatus;
};

export const dropdownInputOptions = (options) => {
  return Object.entries(options).map(([key, value]) => ({ key: key, value: value }));
};

export const filenameFromPath = (filePath) => {
  return path.basename(filePath);
};

// Formats the numeric age input to a human readable string. Eg. 1.5 = 1 day, 2.5 = 2 days
export const formatAgeToDays = (age) => {
  if (age < 1) {
    return 'Less than 1 day';
  }
  if (age >= 1 && age < 2) {
    return '1 day';
  }
  return `${Math.floor(age)} days`;
};

export const formatAddressShort = (address) => {
  const { city, state, postalCode } = address;
  return `${city}, ${state} ${postalCode}`;
};

export const formatPrimeAPIFullAddress = (address) => {
  const { streetAddress1, streetAddress2, city, state, postalCode } = address;
  return `${streetAddress1}, ${streetAddress2}, ${city}, ${state} ${postalCode}`;
};

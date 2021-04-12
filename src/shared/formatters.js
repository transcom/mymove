import { isFinite } from 'lodash';
import moment from 'moment';
import numeral from 'numeral';
import path from 'path';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { DEPARTMENT_INDICATOR_OPTIONS, DEPARTMENT_INDICATOR_LABELS } from 'constants/departmentIndicators';
import { ORDERS_TYPE_OPTIONS, ORDERS_TYPE_DETAILS_OPTIONS } from 'constants/orders';
import { PAYMENT_REQUEST_STATUS_LABELS } from 'constants/paymentRequestStatus';
import { SERVICE_MEMBER_AGENCY_LABELS } from 'content/serviceMemberAgencies';
import { MOVE_STATUS_OPTIONS, SERVICE_COUNSELING_MOVE_STATUS_OPTIONS } from 'constants/queues';

/**
 * Formats number into a dollar string. Eg. $1,234.12
 *
 * More info: http://numeraljs.com/
 * @param num
 * @returns {string}
 */
export function toDollarString(num) {
  return numeral(num).format('$0,0.00');
}

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
    return;
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

// Service Member Formatters

// Format a date in the MM-DD-YYYY format for use in the service member UI.
export function formatDateSM(date) {
  if (date) {
    return moment(date).format('MM/DD/YYYY');
  }
}

// Format a date into the format required for submission as a date property in
// Swagger.
export function formatSwaggerDate(date) {
  if (date) {
    return moment(date).format('YYYY-MM-DD');
  }
  return '';
}

// Parse a date from the format used by Swagger into a Date object
export function parseSwaggerDate(dateString) {
  if (dateString) {
    return moment(dateString, 'YYYY-MM-DD').toDate();
  }
}

// Format a weight with lbs following, e.g. 4000 becomes 4,000 lbs
export function formatWeight(weight) {
  if (weight) {
    return `${weight.toLocaleString()} lbs`;
  } else {
    return '0 lbs';
  }
}

// Format date for display of dates summaries
const formatDateForDateRange = (date, formatType) => {
  let format = '';
  switch (formatType) {
    case 'long':
      format = 'ddd, MMM DD';
      break;
    case 'condensed':
      format = 'MMM D';
      break;
    default:
      format = 'ddd, MMM DD';
  }
  if (date) {
    return moment(date).format(format);
  }
};

export const displayDateRange = (dates, formatType = 'long') => {
  let span = '';
  let firstDate = '';
  if (dates.length > 1) {
    span = ` - ${formatDateForDateRange(dates[dates.length - 1], formatType)}`;
  }
  if (dates.length >= 1) {
    firstDate = formatDateForDateRange(dates[0], formatType);
  }
  return firstDate + span;
};

// Office Formatters

// Format a date and ignore any time values, e.g. 03-Jan-18
export function formatDate(date, inputFormat, outputFormat = 'DD-MMM-YY', locale = 'en', isStrict = false) {
  if (date) {
    return moment(date, inputFormat, locale, isStrict).format(outputFormat);
  }
}

export function formatDateFromIso(date, outputFormat) {
  return formatDate(date, 'YYYY-MM-DDTHH:mm:ss.SSSZ', outputFormat);
}

export function formatDate4DigitYear(date) {
  if (date) {
    return moment(date).format('DD-MMM-YYYY');
  }
}

export function formatTime(date) {
  if (date) {
    return moment(date).format('HH:mm');
  }
}

// Format a date and include its time, e.g. 03-Jan-2018 21:23
export function formatDateTime(date) {
  if (date) {
    return moment(date).format('DD-MMM-YY HH:mm');
  }
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

// maps int to int with ordinal 1 -> 1st, 2 -> 2nd, 3rd ...
export const formatToOrdinal = (n) => {
  const s = ['th', 'st', 'nd', 'rd'];
  const v = n % 100;
  // eslint-disable-next-line security/detect-object-injection
  return n + (s[(v - 20) % 10] || s[v] || s[0]);
};

// Map shipment types to friendly display names for mto shipments
export const mtoShipmentTypeToFriendlyDisplay = (shipmentType) => {
  switch (shipmentType) {
    case SHIPMENT_OPTIONS.HHG:
      return 'Household goods';
    case SHIPMENT_OPTIONS.NTSR:
      return 'NTS release';
    case SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC:
      return 'Household goods longhaul domestic';
    case SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC:
      return 'Household goods shorthaul domestic';
    default:
      return shipmentType;
  }
};

export const departmentIndicatorReadable = (departmentIndicator) => {
  if (!departmentIndicator) {
    return 'Missing';
  }
  return DEPARTMENT_INDICATOR_OPTIONS[`${departmentIndicator}`] || departmentIndicator;
};

export const departmentIndicatorLabel = (departmentIndicator) => {
  return DEPARTMENT_INDICATOR_LABELS[`${departmentIndicator}`] || departmentIndicator;
};

export const serviceMemberAgencyLabel = (agency) => {
  return SERVICE_MEMBER_AGENCY_LABELS[`${agency}`] || agency;
};

export const moveStatusLabel = (status) => {
  return MOVE_STATUS_OPTIONS.find((option) => option.value === `${status}`)?.label || status;
};

export const serviceCounselingMoveStatusLabel = (status) => {
  return SERVICE_COUNSELING_MOVE_STATUS_OPTIONS.find((option) => option.value === `${status}`)?.label || status;
};

export const ordersTypeReadable = (ordersType) => {
  if (!ordersType) {
    return 'Missing';
  }
  return ORDERS_TYPE_OPTIONS[`${ordersType}`] || ordersType;
};

export const ordersTypeDetailReadable = (ordersTypeDetail) => {
  if (!ordersTypeDetail) {
    return 'Missing';
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

export const formatDaysInTransit = (days) => {
  if (days) {
    if (days === 1) {
      return '1 day';
    } else {
      return `${days} days`;
    }
  } else {
    return '0 days';
  }
};

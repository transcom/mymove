import path from 'path';

import moment from 'moment';
import numeral from 'numeral';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { DEPARTMENT_INDICATOR_LABELS, DEPARTMENT_INDICATOR_OPTIONS } from 'constants/departmentIndicators';
import { SERVICE_MEMBER_AGENCY_LABELS } from 'content/serviceMemberAgencies';
import { SERVICE_COUNSELING_MOVE_STATUS_OPTIONS, MOVE_STATUS_OPTIONS } from 'constants/queues';
import { ORDERS_TYPE_OPTIONS, ORDERS_TYPE_DETAILS_OPTIONS } from 'constants/orders';
import { PAYMENT_REQUEST_STATUS_LABELS } from 'constants/paymentRequestStatus';

export function formatNumber(num) {
  if (!Number.isFinite(num)) {
    return '';
  }

  return num.toLocaleString();
}

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

// Format user-entered base quantity into base quantity, e.g. 16.7000 -> 167000
export function formatToBaseQuantity(baseQuantity) {
  return parseFloat(String(baseQuantity).replace(',', '')) * 10000;
}

// Format a base quantity into a user-friendly number, e.g. 167000 -> 16.7
export function convertFromBaseQuantity(baseQuantity) {
  if (!Number.isFinite(baseQuantity)) {
    return null;
  }

  return baseQuantity / 10000;
}

// Format a thousandth of an inch into an inch, e.g. 16700 -> 16.7
export function convertFromThousandthInchToInch(thousandthInch) {
  if (!Number.isFinite(thousandthInch)) {
    return null;
  }

  return thousandthInch / 1000;
}

// Format a dimensions object length, width and height to inches
export function formatToDimensionsInches(dimensions) {
  if (!dimensions) {
    return undefined;
  }

  return {
    length: convertFromThousandthInchToInch(dimensions.length),
    width: convertFromThousandthInchToInch(dimensions.width),
    height: convertFromThousandthInchToInch(dimensions.height),
  };
}

// Format user-entered dimension into base dimension, e.g. 15.25 -> 15250
export function formatToThousandthInches(val) {
  return parseFloat(String(val).replace(',', '')) * 1000;
}

// Format dimensions object length, width and height to base dimensions
export function formatDimensionsToThousandthInches(dimensions) {
  if (!dimensions) {
    return undefined;
  }

  return {
    length: formatToThousandthInches(dimensions.length),
    width: formatToThousandthInches(dimensions.width),
    height: formatToThousandthInches(dimensions.height),
  };
}

// Service Member Formatters

// Format a date in the MM-DD-YYYY format for use in the service member UI.
export function formatDateSM(date) {
  if (date) {
    return moment(date).format('MM/DD/YYYY');
  }
  return undefined;
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
  return undefined;
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
  return undefined;
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
  return undefined;
}

export function formatDateFromIso(date, outputFormat) {
  return formatDate(date, 'YYYY-MM-DDTHH:mm:ss.SSSZ', outputFormat);
}

export function formatDate4DigitYear(date) {
  if (date) {
    return moment(date).format('DD-MMM-YYYY');
  }
  return undefined;
}

export function formatTime(date) {
  if (date) {
    return moment(date).format('HH:mm');
  }
  return undefined;
}

// Format a date and include its time, e.g. 03-Jan-2018 21:23
export function formatDateTime(date) {
  if (date) {
    return moment(date).format('DD-MMM-YY HH:mm');
  }
  return undefined;
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

  return `${moment(date, moment.ISO_8601, true).format('DD-MMM-YY HH:mm')} ${shortZone}`;
}

export function formatTimeAgo(date) {
  if (!date) return undefined;

  return moment(date)
    .fromNow()
    .replace('minute', 'min')
    .replace(/a min\s/, '1 min ');
}

// maps int to int with ordinal 1 -> 1st, 2 -> 2nd, 3rd ...
export const formatToOrdinal = (n) => {
  const s = ['th', 'st', 'nd', 'rd'];
  const v = n % 100;
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

export const departmentIndicatorReadable = (departmentIndicator, missingText) => {
  if (!departmentIndicator) {
    return missingText;
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

export const filenameFromPath = (filePath) => {
  return path.basename(filePath);
};

export const formatAddressShort = (address) => {
  const { city, state, postalCode } = address;
  return `${city}, ${state} ${postalCode}`;
};

export const formatPrimeAPIFullAddress = (address) => {
  const { streetAddress1, streetAddress2, city, state, postalCode } = address;
  return `${streetAddress1}, ${streetAddress2}, ${city}, ${state} ${postalCode}`;
};

export const formatMoveHistoryFullAddress = (address) => {
  let formattedAddress = '';
  if (address.street_address_1) {
    formattedAddress += `${address.street_address_1}`;
  }

  if (address.street_address_2) {
    formattedAddress += `, ${address.street_address_2}`;
  }

  if (address.city) {
    formattedAddress += `, ${address.city}`;
  }

  if (address.state) {
    formattedAddress += `, ${address.state}`;
  }

  if (address.postal_code) {
    formattedAddress += ` ${address.postal_code}`;
  }

  if (formattedAddress[0] === ',') {
    formattedAddress = formattedAddress.substring(1);
  }

  formattedAddress = formattedAddress.trim();

  return formattedAddress;
};

export const dropdownInputOptions = (options) => {
  return Object.entries(options).map(([key, value]) => ({ key, value }));
};

// adds commas to numberString w/o removeing .0000 from the end of the string or rounding
export const addCommasToNumberString = (numOrString, decimalPlaces = 0) => {
  let numToFixed = numOrString;
  if (!numOrString || numOrString === '0') {
    numToFixed = (0).toFixed(decimalPlaces);
  }

  const str = numToFixed.toString();
  const [wholeNum, decimalNum] = str.split('.');
  const wholeNumInt = parseInt(wholeNum, 10);
  if (decimalNum) {
    return `${wholeNumInt.toLocaleString()}.${decimalNum}`;
  }
  return wholeNumInt.toLocaleString();
};

// truncate a number and return appropiate decimal places... (watch out for negitive numbers: floor(-5.1) === -6)
// see test for examples of how this works
export const truncateNumber = (num, decimalPlaces = 0) => {
  if (!num) return num;

  const floatNum = parseFloat(num).toFixed(4);
  const scale = 10 ** decimalPlaces;
  const truncatedNbr = Math.floor(floatNum * scale) / scale;
  return truncatedNbr.toFixed(decimalPlaces).toString();
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

// Format orders type (ex: PERMANENT_CHANGE_OF_STATION => Permanent change of station)
export function formatOrderType(orderType) {
  return orderType
    .split('_')
    .map((str, i) => {
      if (i === 0) {
        return str[0] + str.slice(1).toLowerCase();
      }
      return str.toLowerCase();
    })
    .join(' ');
}

// Format dates for customer app (ex. 25 Dec 2020)
export function formatCustomerDate(date) {
  return moment(date).format('DD MMM YYYY');
}

export function formatSignatureDate(date) {
  return moment(date).format('YYYY-MM-DD');
}

// Translate boolean (true/false) into "yes"/"no" string
export const formatYesNoInputValue = (value) => {
  if (value === true) return 'yes';
  if (value === false) return 'no';
  return null;
};

// Translate "yes"/"no" string into boolean (true/false)
export const formatYesNoAPIValue = (value) => {
  if (value === 'yes') return true;
  if (value === 'no') return false;
  return undefined;
};

// Translate weights from lbs to CWT
export const formatWeightCWTFromLbs = (value) => {
  return `${parseInt(value, 10) / 100} cwt`;
};

// Translate currency from millicents to dollars
export const formatDollarFromMillicents = (value) => {
  return `$${(parseInt(value, 10) / 100000).toFixed(2)}`;
};

// Takes an whole number of day value and pluralizes with unit label
export const formatDaysInTransit = (days) => {
  if (days) {
    if (days === 1) {
      return '1 day';
    }
    return `${days} days`;
  }
  return '0 days';
};

export const formatAddressForPrimeAPI = (address) => {
  return {
    streetAddress1: address.streetAddress1,
    streetAddress2: address.streetAddress2,
    streetAddress3: address.streetAddress3,
    city: address.city,
    state: address.state,
    postalCode: address.postalCode,
  };
};

const emptyAddress = {
  streetAddress1: '',
  streetAddress2: '',
  city: '',
  state: '',
  postalCode: '',
};

export function fromPrimeAPIAddressFormat(address) {
  if (!address) {
    return emptyAddress;
  }
  return {
    streetAddress1: address.streetAddress1,
    streetAddress2: address.streetAddress2,
    streetAddress3: address.streetAddress3,
    city: address.city,
    state: address.state,
    postalCode: address.postalCode,
  };
}

// Format a weight with lbs following, e.g. 4000 becomes 4,000 lbs
export function formatWeight(weight) {
  if (weight) {
    return `${weight.toLocaleString()} lbs`;
  }
  return '0 lbs';
}

export const formatDelimitedNumber = (number) => {
  // Fail-safe in case an actual number value is passed in
  const numberString = number.toString();
  return Number(numberString.replace(/,/g, ''));
};
/**
 * Depending on the order type, this will return:
 * Report by date (PERMANENT_CHANGE_OF_STATION)
 * Date of retirement (RETIREMENT)
 * Date of separation (SEPARATION)
 */
export const formatLabelReportByDate = (orderType) => {
  switch (orderType) {
    case 'RETIREMENT':
      return 'Date of retirement';
    case 'SEPARATION':
      return 'Date of separation';
    default:
      return 'Report by date';
  }
};

// Format a number of cents into a string, e.g. 12,345.67
export function formatCents(cents, minimumFractionDigits = 2, maximumFractionDigits = 2) {
  return (cents / 100).toLocaleString(undefined, { minimumFractionDigits, maximumFractionDigits });
}

// Format base quantity as cents
export function formatBaseQuantityAsDollars(baseQuantity) {
  return formatCents(baseQuantity / 100);
}

export function formatCentsRange(min, max) {
  if (!Number.isFinite(min) || !Number.isFinite(max)) {
    return '';
  }

  return `$${formatCents(min)} - ${formatCents(max)}`;
}

// Format a base quantity into a user-friendly number string, e.g. 167000 -> "16.7000"
export function formatFromBaseQuantity(baseQuantity) {
  if (!Number.isFinite(baseQuantity)) {
    return '';
  }

  return (baseQuantity / 10000).toLocaleString(undefined, {
    minimumFractionDigits: 4,
    maximumFractionDigits: 4,
  });
}

// Formats a numeric value amount in the default locale with configurable options
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Number/toLocaleString
export function formatAmount(amount, options = { minimumFractionDigits: 2, maximumFractionDigits: 2 }) {
  if (!Number.isFinite(amount)) {
    return '';
  }
  return amount.toLocaleString(undefined, options);
}

// Converts a cents value into whole dollars, dropping the decimal precision without rounding e.g. 1234599 -> 12,345
export function formatCentsTruncateWhole(cents) {
  return formatAmount(Math.floor(cents / 100), { minimumFractionDigits: 0, maximumFractionDigits: 0 });
}

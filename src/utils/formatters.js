import path from 'path';

import moment from 'moment';
import numeral from 'numeral';

import { DEPARTMENT_INDICATOR_OPTIONS } from 'constants/departmentIndicators';
import { SERVICE_MEMBER_AGENCY_LABELS } from 'content/serviceMemberAgencies';
import { ORDERS_TYPE_OPTIONS, ORDERS_TYPE_DETAILS_OPTIONS } from 'constants/orders';
import { PAYMENT_REQUEST_STATUS_LABELS } from 'constants/paymentRequestStatus';

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

// Format a thousandth of an inch into an inch, e.g. 16700 -> 16.7
export function convertFromThousandthInchToInch(thousandthInch) {
  if (!Number.isFinite(thousandthInch)) {
    return null;
  }

  return thousandthInch / 1000;
}

// Format user-entered dimension into base dimension, e.g. 15.25 -> 15250
export function formatToThousandthInches(val) {
  return parseFloat(String(val).replace(',', '')) * 1000;
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

// Format a date and include its time, e.g. 03-Jan-2018 21:23
export function formatDateTime(date) {
  if (date) {
    return moment(date).format('DD-MMM-YY HH:mm');
  }
  return undefined;
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

export const departmentIndicatorReadable = (departmentIndicator, missingText) => {
  if (!departmentIndicator) {
    return missingText;
  }
  return DEPARTMENT_INDICATOR_OPTIONS[`${departmentIndicator}`] || departmentIndicator;
};

export const serviceMemberAgencyLabel = (agency) => {
  return SERVICE_MEMBER_AGENCY_LABELS[`${agency}`] || agency;
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

export const formatMoveHistoryAgent = (agent) => {
  let formattedAgent = '';

  if (agent.first_name) {
    formattedAgent += `${agent.first_name}`;
  }

  if (agent.last_name) {
    formattedAgent += ` ${agent.last_name}`;
  }

  if (agent.phone) {
    formattedAgent += `, ${agent.phone}`;
  }

  if (agent.email) {
    formattedAgent += `, ${agent.email}`;
  }

  if (formattedAgent[0] === ',') {
    formattedAgent = formattedAgent.substring(1);
  }

  formattedAgent = formattedAgent.trim();

  return formattedAgent;
};

export const dropdownInputOptions = (options) => {
  return Object.entries(options).map(([key, value]) => ({ key, value }));
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

// Format dates for customer app (ex. 25 Dec 2020)
export function formatCustomerDate(date) {
  return moment(date).format('DD MMM YYYY');
}
// Format dates for customer remarks in the office app (ex. 25 Dec 2020 8:00)
export function formatCustomerSupportRemarksDate(date) {
  return moment(date).format('DD MMM YYYY HH:mm');
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

export function formatCentsRange(min, max) {
  if (!Number.isFinite(min) || !Number.isFinite(max)) {
    return '';
  }

  return `$${formatCents(min)} - ${formatCents(max)}`;
}

// Formats a numeric value amount in the default locale with configurable options
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Number/toLocaleString
export function formatAmount(amount, options = { minimumFractionDigits: 2, maximumFractionDigits: 2 }) {
  if (!Number.isFinite(amount)) {
    return '';
  }
  return amount.toLocaleString(undefined, options);
}

// Converts a cents value into whole dollars, rounding down.
export function convertCentsToWholeDollarsRoundedDown(cents) {
  return Math.floor(cents / 100);
}

// Converts a cents value into whole dollars, dropping the decimal precision without rounding e.g. 1234599 -> 12,345
export function formatCentsTruncateWhole(cents) {
  return formatAmount(convertCentsToWholeDollarsRoundedDown(cents), {
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  });
}

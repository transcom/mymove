import moment from 'moment';
import numeral from 'numeral';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { DEPARTMENT_INDICATOR_LABELS, DEPARTMENT_INDICATOR_OPTIONS } from 'constants/departmentIndicators';
import { SERVICE_MEMBER_AGENCY_LABELS } from 'content/serviceMemberAgencies';
import { SERVICE_COUNSELING_MOVE_STATUS_OPTIONS, MOVE_STATUS_OPTIONS } from 'constants/queues';
import { ORDERS_TYPE_OPTIONS } from 'constants/orders';

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

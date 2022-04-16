import moment from 'moment';
import numeral from 'numeral';

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

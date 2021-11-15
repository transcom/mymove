import moment from 'moment';

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
    streetAddress1: address.street_address_1,
    streetAddress2: address.street_address_2,
    streetAddress3: address.street_address_3,
    city: address.city,
    state: address.state,
    postalCode: address.postal_code,
  };
};

const emptyAddress = {
  street_address_1: '',
  street_address_2: '',
  city: '',
  state: '',
  postal_code: '',
};

export function fromPrimeAPIAddressFormat(address) {
  if (!address) {
    return emptyAddress;
  }
  return {
    street_address_1: address.streetAddress1,
    street_address_2: address.streetAddress2,
    street_address_3: address.streetAddress3,
    city: address.city,
    state: address.state,
    postal_code: address.postalCode,
  };
}

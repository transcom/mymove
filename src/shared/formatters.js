import { isFinite } from 'lodash';
import moment from 'moment';

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
      format = 'MMM DD';
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
export function formatDate(date) {
  if (date) {
    return moment(date).format('DD-MMM-YY');
  }
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

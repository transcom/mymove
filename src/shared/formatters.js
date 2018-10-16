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

// Format a base quantity into a user-friendly number, e.g. 167000 -> 16.7000
export function formatFromBaseQuantity(baseQuantity) {
  if (!isFinite(baseQuantity)) {
    return '';
  }

  return (baseQuantity / 10000).toLocaleString(undefined, {
    minimumFractionDigits: 4,
    maximumFractionDigits: 4,
  });
}

// Format user-entered base quantity into base quantity, e.g. 16.7000 -> 167000
export function formatToBaseQuantity(baseQuantity) {
  baseQuantity = parseFloat(String(baseQuantity).replace(',', '')) * 10000;

  return baseQuantity;
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

// Office Formatters

// Format a date and ignore any time values, e.g. 03-Jan-2018
export function formatDate(date) {
  if (date) {
    return moment(date).format('DD-MMM-YY');
  }
}

// Format a date and include its time, e.g. 03-Jan-2018 21:23
export function formatDateTime(date) {
  if (date) {
    return moment(date).format('DD-MMM-YY HH:mm');
  }
}

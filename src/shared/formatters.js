import moment from 'moment';

// Format a number of cents into a string, e.g. $12,345.67
export function formatCents(cents) {
  if (!cents) {
    return '';
  }

  return (cents / 100).toLocaleString(undefined, {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });
}

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

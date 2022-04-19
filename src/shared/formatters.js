import { isFinite } from 'lodash';
import moment from 'moment';

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

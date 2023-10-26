import moment from 'moment';

export const swaggerDateFormat = 'YYYY-MM-DD';
export const defaultDateFormat = 'M/D/YYYY';
export const utcDateFormat = 'YYYY-MM-DDTHH:mm:ssZ';
export const datePickerFormat = 'DD MMM YYYY';

// First date format is take to be the default
const allowedDateFormats = [
  defaultDateFormat,
  'YYYY/M/D',
  'YYYY-M-D',
  'M-D-YYYY',
  'D-MMM-YYYY',
  'MMM-D-YYYY',
  'DD-MMM-YY',
  'DD MMM YYYY',
  'YYYY-MM-DDTHH:mm:ssZ',
  'YYYY-MM-DDTHH:mm:ss.SSSZ',
];

export function parseDate(str, _format, locale = 'en') {
  // Ignore default format, and attempt to parse date using allowed formats
  const m = moment(str, allowedDateFormats, locale, true);
  if (m.isValid()) {
    return m.toDate();
  }

  return undefined;
}

export function formatDate(date, format = defaultDateFormat, locale = 'en') {
  return moment(date, allowedDateFormats, locale, true).locale(locale).format(format);
}

export function formatDateWithUTC(date, format = defaultDateFormat, locale = 'en') {
  return moment.utc(date, allowedDateFormats, locale, true).locale(locale).format(format);
}

export function formatDateForSwagger(dateString) {
  if (dateString) {
    return formatDate(dateString, swaggerDateFormat);
  }
  return undefined;
}

export function formatDateTime(dateString) {
  if (dateString) {
    const startOfDay = moment(new Date(dateString)).hour(0);
    return moment.utc(startOfDay).format(utcDateFormat);
  }
  return undefined;
}

/**
 * @function
 * @description This function is to convert dates to strings in the format used
 * by the DatePickerInput
 * @param {moment.input} date A Moment.input representing a date
 * @returns {String} A String representing the date in the string format used by the DatePickerInput
 */
export function formatDateForDatePicker(date) {
  if (date) {
    return formatDate(date, datePickerFormat);
  }
  return undefined;
}

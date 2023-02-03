import moment from 'moment';

export const swaggerDateFormat = 'YYYY-MM-DD';
export const defaultDateFormat = 'M/D/YYYY';
export const utcDateFormat = 'YYYY-MM-DDTHH:mm:ssZ';

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

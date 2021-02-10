import dayjs from 'dayjs';
import customParseFormat from 'dayjs/plugin/customParseFormat';

dayjs.extend(customParseFormat);

export const swaggerDateFormat = 'YYYY-MM-DD';
export const defaultDateFormat = 'M/D/YYYY';
export const ISO_8601_FORMAT = 'YYYY-MM-DDTHH:mm:ss.SSSZ';

// First date format is take to be the default
const allowedDateFormats = [
  defaultDateFormat,
  ISO_8601_FORMAT,
  'YYYY/M/D',
  'YYYY-M-D',
  'YYYY-MM-DD',
  'M-D-YYYY',
  'D-MMM-YYYY',
  'MMM-D-YYYY',
  'DD-MMM-YY',
  'DD MMM YYYY',
];

export function parseDate(str, _format, locale = 'en') {
  // Ignore _format parameter, and attempt to parse date using allowed formats
  const m = dayjs(str, allowedDateFormats, locale);

  if (m.isValid()) {
    return m.toDate();
  }

  return undefined;
}

export function formatDateFromISO(date, outputFormat = defaultDateFormat, locale = 'en') {
  return dayjs(date).locale(locale).format(outputFormat);
}

export function formatDate(date, format = defaultDateFormat, locale = 'en') {
  return dayjs(date, allowedDateFormats, locale, true).locale(locale).format(format);
}

export function formatDateForSwagger(dateString) {
  if (dateString) {
    return formatDate(dateString, swaggerDateFormat);
  }
}

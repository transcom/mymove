/*
  When you need to put some control between what the user enters and the value that gets stored in
  Redux, you can use a "normalizer". A normalizer is just a function that gets run every time a value
  is changed that can transform the value before storing.

  For more information: https://redux-form.com/7.4.2/examples/normalizing/
*/

import { swaggerDateFormat } from 'shared/utils';
import moment from 'moment';

const normalizePhone = value => {
  if (!value) {
    return value;
  }
  const onlyNums = value.replace(/[^\d]/g, '');
  let normalizedPhone = '';
  for (let i = 0; i < 10; i++) {
    if (i >= onlyNums.length) {
      break;
    }
    if (i === 3 || i === 6) {
      normalizedPhone += '-';
    }
    normalizedPhone += onlyNums[i]; // eslint-disable-line security/detect-object-injection
  }
  return normalizedPhone;
};

const normalizeSSN = value => {
  if (!value) {
    return value;
  }
  const onlyNums = value.replace(/[^\d]/g, '');
  let normalizedSSN = '';
  for (let i = 0; i < 9; i++) {
    if (i >= onlyNums.length) {
      break;
    }
    if (i === 3 || i === 5) {
      normalizedSSN += '-';
    }
    normalizedSSN += onlyNums[i]; // eslint-disable-line security/detect-object-injection
  }
  return normalizedSSN;
};

const normalizeZip = value => {
  if (!value) {
    return value;
  }
  const onlyNums = value.replace(/[^\d]/g, '');
  let normalizedZip = '';
  for (let i = 0; i < 9; i++) {
    if (i >= onlyNums.length) {
      break;
    }
    if (i === 5) {
      normalizedZip += '-';
    }
    normalizedZip += onlyNums[i]; // eslint-disable-line security/detect-object-injection
  }
  return normalizedZip;
};

const normalizeDates = value => {
  return value ? moment(value).format(swaggerDateFormat) : value;
};

const createDigitNormalizer = maxLength => {
  return value => {
    if (!value) {
      return value;
    }

    // only digits up to the max length
    // if undefined, max length is length of string
    return value.replace(/[^\d]/g, '').substr(0, maxLength);
  };
};

const createDecimalNormalizer = decimalDigits => {
  return value => {
    if (!value) {
      return value;
    }
    value = value.replace(/[^\d.]/g, '');

    if (value.indexOf('.') >= 0) {
      value =
        value.substr(0, value.indexOf('.')) +
        '.' +
        value
          .substr(value.indexOf('.') + 1)
          .replace(/[^\d]/g, '')
          .substr(0, decimalDigits);
    }
    return value;
  };
};

export { normalizePhone, normalizeSSN, normalizeZip, normalizeDates, createDecimalNormalizer, createDigitNormalizer };

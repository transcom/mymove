import { memoize } from 'lodash';

const isRequired = value => (value ? undefined : 'Required');
// Why Memoize? Please see https://github.com/erikras/redux-form/issues/3288
// Since we attach validators inside the render method, without memoization the
// function is re-created on every render which is not handled by react form.
// By memoizing it, it works.
const maxLength = memoize(maxLength => value => {
  if (value && value.length > maxLength) {
    return `Cannot exceed ${maxLength} characters.`;
  }
});
const minLength = memoize(minLength => value => {
  if (value && value.length < minLength) {
    return `Must be at least ${minLength} characters long.`;
  }
});

const maximum = memoize(maximum => value => {
  if (value && value > maximum) {
    return `Must be ${maximum} or less`;
  }
});
const minimum = memoize(minimum => value => {
  if (value && value < minimum) {
    return `Must be ${minimum} or more`;
  }
});

const isNumber = value => {
  if (value) {
    if (isNaN(parseFloat(value))) {
      return 'Must be a number.';
    }
  }
};

const isInteger = value => {
  if (value) {
    if (!Number.isInteger(value)) {
      return 'Must be an integer';
    }
  }
};

const normalizePhone = (value, previousValue) => {
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

const normalizeSSN = (value, previousValue) => {
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

const normalizeZip = (value, previousValue) => {
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

const patternMatches = memoize((pattern, message) => {
  const regex = RegExp(pattern);
  return value => {
    if (value) {
      if (!regex.test(value)) {
        return message;
      }
    }
  };
});

export default {
  maxLength,
  minLength,
  maximum,
  minimum,
  isRequired,
  isNumber,
  isInteger,
  normalizePhone,
  normalizeSSN,
  normalizeZip,
  patternMatches,
};

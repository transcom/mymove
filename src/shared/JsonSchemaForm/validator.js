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

const isPhoneNumber = value => {
  if (value && value.replace(/[^\d]/g, '').length !== 10) {
    return 'Number must have 10 digits.';
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
    normalizedPhone += onlyNums[i];
  }
  return normalizedPhone;
};

const isSSN = value => {
  if (value && value.replace(/[^\d]/g, '').length !== 9) {
    return 'SSN must have 9 digits.';
  }
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
    normalizedSSN += onlyNums[i];
  }
  return normalizedSSN;
};

const isZip = value => {
  if (value) {
    const zipLength = value.replace(/[^\d]/g, '').length;
    if (!(zipLength === 9 || zipLength === 5)) {
      return 'Zip code must have 5 or 9 digits.';
    }
  }
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
    normalizedZip += onlyNums[i];
  }
  return normalizedZip;
};

const patternMatches = memoize((pattern, example) => {
  console.log('patternavlaid');
  const regex = RegExp(pattern);
  return value => {
    console.log('patternavlaidsssss');
    if (!regex.test(value)) {
      return 'Incorrect Format. This is a valid example: ' + example;
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
  isPhoneNumber,
  normalizePhone,
  isSSN,
  normalizeSSN,
  isZip,
  normalizeZip,
  patternMatches,
};

import { isFinite, isInteger as rawIsInteger, memoize } from 'lodash';
import { defaultDateFormat } from 'shared/dates';
import moment from 'moment';

const isRequired = (value) => (value ? undefined : 'Required');
// Why Memoize? Please see https://github.com/erikras/redux-form/issues/3288
// Since we attach validators inside the render method, without memoization the
// function is re-created on every render which is not handled by react form.
// By memoizing it, it works.
const maxLength = memoize((maxLength) => (value) => {
  if (value && value.length > maxLength) {
    return `Cannot exceed ${maxLength} characters.`;
  }
  return undefined;
});
const minLength = memoize((minLength) => (value) => {
  if (value && value.length < minLength) {
    return `Must be at least ${minLength} characters long.`;
  }
  return undefined;
});

const maximum = memoize((maximum) => (value) => {
  if (value && value > maximum) {
    return `Must be ${maximum} or less`;
  }
  return undefined;
});
const minimum = memoize((minimum) => (value) => {
  if (value && value < minimum) {
    return `Must be ${minimum} or more`;
  }
  return undefined;
});

const isNumber = (value) => {
  if (value) {
    if (!isFinite(parseFloat(value))) {
      return 'Must be a number';
    }
  }
  return undefined;
};

const isInteger = (value) => {
  if (value) {
    if (!rawIsInteger(value)) {
      return 'Must be an integer';
    }
  }
  return undefined;
};

const isDate = (value) => {
  if (value) {
    let parsed = moment(value, defaultDateFormat);
    if (!parsed.isValid()) {
      return 'Must be a valid date';
    }
  }
  return undefined;
};

const patternMatches = memoize((pattern, message) => {
  const regex = RegExp(pattern);
  return (value) => {
    if (value) {
      if (!regex.test(value)) {
        return message;
      }
    }
    return undefined;
  };
});

const minDateValidation = memoize((minDate = null, message) => {
  return (value) => {
    if (minDate && moment(value).isBefore(moment(minDate))) {
      return message;
    }
    return undefined;
  };
});

const maxDateValidation = memoize((maxDate = null, message) => {
  return (value) => {
    if (maxDate && moment(value).isAfter(moment(maxDate))) {
      return message;
    }
    return undefined;
  };
});

const validators = {
  maxLength,
  minLength,
  maximum,
  minimum,
  isRequired,
  isNumber,
  isInteger,
  isDate,
  minDateValidation,
  maxDateValidation,
  patternMatches,
};

export default validators;

const INVALID_DATE = 'Invalid date';

export function validateDate(value) {
  let error;
  if (value === INVALID_DATE || !value) {
    error = 'Required';
  }
  return error;
}

export function validateNotEmpty(value) {
  let error;
  if (!value) {
    error = 'Required';
  }
  return error;
}

export function validateState(value) {
  let error;
  if (!value) {
    error = 'Required';
  } else if (value.length !== 2) {
    error = 'Must be state abbreviation';
  }
  return error;
}

export function validateZIPCode(value) {
  // eslint-disable-next-line security/detect-unsafe-regex
  const validZipCode = RegExp(/^(\d{5}([-]\d{4})?)$/);
  let error;
  if (!value) {
    error = 'Required';
  } else if (!validZipCode.test(value)) {
    error = 'Must be valid ZIP code';
  }
  return error;
}

const INVALID_DATE = 'Invalid date';

// eslint-disable-next-line import/prefer-default-export
export function validateDate(value) {
  let error;
  if (value === INVALID_DATE || !value) {
    error = 'Required';
  }
  return error;
}

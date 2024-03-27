import { required, regex, minLength, maxLength } from 'react-admin';

export const phoneValidators = [
  required(),
  regex(/^[2-9]\d{2}-\d{3}-\d{4}$/, 'Invalid phone number, should be 000-000-0000'),
];

export const edipiValidator = [minLength('10'), maxLength('10')];

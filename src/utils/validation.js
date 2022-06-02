import * as Yup from 'yup';

import { ValidateZipRateData } from 'shared/api';

const INVALID_DATE = 'Invalid date';

// RA Summary: eslint - security/detect-unsafe-regex - Denial of Service: Regular Expression
// RA: Locates potentially unsafe regular expressions, which may take a very long time to run, blocking the event loop
// RA: Per MilMove SSP, predisposing conditions are regex patterns from untrusted sources or unbounded matching.
// RA: The regex pattern is a constant string set at compile-time and it is bounded to 10 characters (zip code).
// RA Developer Status: Mitigated
// RA Validator Status:  Mitigated
// RA Modified Severity: N/A
// eslint-disable-next-line security/detect-unsafe-regex
export const ZIP_CODE_REGEX = /^(\d{5}([-]\d{4})?)$/;

export const ZIP5_CODE_REGEX = /^(\d{5})$/;

// eslint-disable-next-line import/prefer-default-export
export function validateDate(value) {
  let error;
  if (value === INVALID_DATE || !value) {
    error = 'Required';
  }
  return error;
}

export const UnsupportedZipCodeErrorMsg =
  'Sorry, we don’t support that zip code yet. Please contact your local PPPO for assistance.';

export const UnsupportedZipCodePPMErrorMsg =
  "We don't have rates for this ZIP code. Please verify that you have entered the correct one. Contact support if this problem persists.";

export const InvalidZIPTypeError = 'Enter a 5-digit ZIP code';

export const validatePostalCode = async (value, postalCodeType, errMsg = UnsupportedZipCodeErrorMsg) => {
  if (!value || (value.length !== 5 && value.length !== 10)) {
    return undefined;
  }

  let responseBody;
  try {
    responseBody = await ValidateZipRateData(value, postalCodeType);
  } catch (e) {
    return 'Error checking ZIP';
  }

  return responseBody.valid ? undefined : errMsg;
};

/** Yup validation schemas */

export const requiredAddressSchema = Yup.object().shape({
  streetAddress1: Yup.string().required('Required'),
  streetAddress2: Yup.string(),
  city: Yup.string().required('Required'),
  state: Yup.string().length(2, 'Must use state abbreviation').required('Required'),
  postalCode: Yup.string().matches(ZIP_CODE_REGEX, 'Must be valid zip code').required('Required'),
});

export const addressSchema = Yup.object().shape({
  streetAddress1: Yup.string(),
  streetAddress2: Yup.string(),
  city: Yup.string(),
  state: Yup.string().length(2, 'Must use state abbreviation'),
  postalCode: Yup.string().matches(ZIP_CODE_REGEX, 'Must be valid zip code'),
});

export const phoneSchema = Yup.string().min(12, 'Number must have 10 digits and a valid area code'); // min 12 includes hyphens

export const emailSchema = Yup.string().matches(
  /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$/,
  'Must be a valid email address',
);

const validatePreferredContactMethod = (value, testContext) => {
  return testContext.parent.phone_is_preferred || testContext.parent.email_is_preferred;
};

export const contactInfoSchema = Yup.object().shape({
  telephone: phoneSchema.required('Required'),
  secondary_telephone: phoneSchema,
  personal_email: emailSchema.required('Required'),
  phone_is_preferred: Yup.bool().test(
    'contactMethodRequired',
    'Please select a preferred method of contact.',
    validatePreferredContactMethod,
  ),
  email_is_preferred: Yup.bool().test(
    'contactMethodRequired',
    'Please select a preferred method of contact.',
    validatePreferredContactMethod,
  ),
});

export const backupContactInfoSchema = Yup.object().shape({
  name: Yup.string().required('Required'),
  email: emailSchema.required('Required'),
  telephone: phoneSchema.required('Required'),
});

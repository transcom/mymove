import * as Yup from 'yup';

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

// eslint-disable-next-line import/prefer-default-export
export function validateDate(value) {
  let error;
  if (value === INVALID_DATE || !value) {
    error = 'Required';
  }
  return error;
}

/** Yup validation schemas */

export const requiredAddressSchema = Yup.object().shape({
  street_address_1: Yup.string().required('Required'),
  street_address_2: Yup.string(),
  city: Yup.string().required('Required'),
  state: Yup.string().length(2, 'Must use state abbreviation').required('Required'),
  postal_code: Yup.string().matches(ZIP_CODE_REGEX, 'Must be valid zip code').required('Required'),
});

export const addressSchema = Yup.object().shape({
  street_address_1: Yup.string(),
  street_address_2: Yup.string(),
  city: Yup.string(),
  state: Yup.string().length(2, 'Must use state abbreviation'),
  postal_code: Yup.string().matches(ZIP_CODE_REGEX, 'Must be valid zip code'),
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

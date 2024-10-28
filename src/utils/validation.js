import * as Yup from 'yup';

import { getUnSupportedStates, unSupportedStates, unSupportedStatesDisabledAlaska } from '../constants/states';

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

export const PHONE_NUMBER_REGEX = /^[2-9]\d{2}-\d{3}-\d{4}$/;

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

export const UnsupportedStateErrorMsg = 'Moves to this state are not supported at this time.';
export const IsSupportedState = async (value, context) => {
  const selectedState = value;

  const enableAK = 'enabledAK';

  let unsupportedStates;
  if (enableAK in context.options.context) {
    unsupportedStates = context.options.context.enabledAK ? unSupportedStates : unSupportedStatesDisabledAlaska;
  } else {
    unsupportedStates = await getUnSupportedStates();
  }

  const found = unsupportedStates.find((unsupportedState) => unsupportedState.key === selectedState);

  if (found) {
    return false;
  }

  return true;
};

/** Yup validation schemas */

export const requiredAddressSchema = Yup.object().shape({
  streetAddress1: Yup.string().trim().required('Required'),
  streetAddress2: Yup.string(),
  city: Yup.string().trim().required('Required'),
  state: Yup.string()
    .test('', UnsupportedStateErrorMsg, IsSupportedState)
    .length(2, 'Must use state abbreviation')
    .required('Required'),
  postalCode: Yup.string().matches(ZIP_CODE_REGEX, 'Must be valid zip code').required('Required'),
});

export const requiredW2AddressSchema = Yup.object().shape({
  streetAddress1: Yup.string().required('Required'),
  streetAddress2: Yup.string(),
  city: Yup.string().required('Required'),
  state: Yup.string().length(2, 'Must use state abbreviation').required('Required'),
  postalCode: Yup.string().matches(ZIP5_CODE_REGEX, 'Must be valid zip code').required('Required'),
});

export const addressSchema = Yup.object().shape({
  streetAddress1: Yup.string(),
  streetAddress2: Yup.string(),
  city: Yup.string(),
  state: Yup.string().length(2, 'Must use state abbreviation'),
  postalCode: Yup.string().matches(ZIP_CODE_REGEX, 'Must be valid zip code'),
});

export const phoneSchema = Yup.string().matches(
  PHONE_NUMBER_REGEX,
  'Please enter a valid phone number. Phone numbers must be entered as ###-###-####.',
); // min 12 includes hyphens

export const OfficeAccountRequestEmailSchema = Yup.string().matches(
  /^[a-zA-Z0-9._%+-]+@(.[a-zA-Z0-9-.]+)[.]{1}(?<!gov|edu|mil)(gov|edu|mil)(?!gov|edu|mil)$/,
  'Domain must be .mil, .gov or .edu',
);

export const emailSchema = Yup.string().matches(
  /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$/,
  'Must be a valid email address',
);

export const preferredContactMethodValidation = (value) => {
  return (
    value?.phone_is_preferred ||
    value?.email_is_preferred ||
    new Yup.ValidationError('Please select a preferred method of contact.', null, 'preferredContactMethod')
  );
};

export const contactInfoSchema = Yup.object()
  .shape({
    telephone: phoneSchema.required('Required'),
    secondary_telephone: phoneSchema,
    personal_email: emailSchema.required('Required'),
    phone_is_preferred: Yup.bool(),
    email_is_preferred: Yup.bool(),
  })
  .test('contactMethodRequired', 'Please select a preferred method of contact.', preferredContactMethodValidation);

export const backupContactInfoSchema = Yup.object().shape({
  name: Yup.string().required('Required'),
  email: emailSchema.required('Required'),
  telephone: phoneSchema.required('Required'),
});

export const edipiMaxErrorMsg = 'Must be 10 digits in length';
export const emailFormatErrorMsg = 'Must be in email format';
export const numericOnlyErrorMsg = 'EDIPI must contain only numeric characters';
export const noNumericAllowedErrorMsg = 'Cannot contain numeric characters';
export const domainFormatErrorMsg = 'Email address must end in a valid domain';
export const allowedDomains = ['.com', '.gov', '.mil', '.edu', '.org', '.net', '.int', '.eu', '.io', '.co'];

// checking okta profile edit form
// oktaEmail must end in the domain listed in allowedDomain variable
// oktaFirst&LastName must not contain numbers
// edipi can only be numbers
// we are validating here to avoid confusing swagger errors
export const oktaInfoSchema = Yup.object().shape({
  oktaUsername: Yup.string().required('Required'),
  oktaEmail: Yup.string()
    .test('domain-suffix', domainFormatErrorMsg, (value) => {
      if (!value) {
        return true;
      }
      const domainMatch = value.match(/@([A-Za-z0-9.-]+)$/);

      if (domainMatch) {
        const domain = domainMatch[1].toLowerCase();
        const tldMatch = domain.match(/\.[A-Za-z]+$/);
        if (tldMatch) {
          const tld = tldMatch[0].toLowerCase();
          return allowedDomains.includes(tld);
        }
      }
      return false;
    })
    .email(emailFormatErrorMsg)
    .required('Required'),
  oktaFirstName: Yup.string()
    .matches(/^[A-Za-z]+$/, noNumericAllowedErrorMsg)
    .required('Required'),
  oktaLastName: Yup.string()
    .matches(/^[A-Za-z]+$/, noNumericAllowedErrorMsg)
    .required('Required'),
  oktaEdipi: Yup.string()
    .min(10, edipiMaxErrorMsg)
    .max(10, edipiMaxErrorMsg)
    .matches(/^[0-9]*$/, numericOnlyErrorMsg)
    .nullable(),
});

export const otherUniqueIdErrorMsg = 'Only accepts alphanumeric characters';
export const middleInitialErrorMsg = 'Must be a single uppercase character';

// Validation method for Office Account Request Form checkbox fields
const validateRoleRequestedMethod = (value, testContext) => {
  return (
    testContext.parent.taskOrderingOfficerCheckBox ||
    testContext.parent.taskInvoicingOfficerCheckBox ||
    testContext.parent.servicesCounselorCheckBox ||
    testContext.parent.transportationContractingOfficerCheckBox ||
    testContext.parent.qualityAssuranceEvaluatorCheckBox ||
    testContext.parent.headquartersCheckBox ||
    testContext.parent.customerSupportRepresentativeCheckBox ||
    testContext.parent.governmentSurveillanceRepresentativeCheckbox
  );
};

const validateOtherUniqueId = (value, testContext) => {
  if (testContext.parent.officeAccountRequestOtherUniqueId || testContext.parent.officeAccountRequestEdipi) {
    return true;
  }

  return false;
};

const validateEdipi = (value, testContext) => {
  if (testContext.parent.officeAccountRequestOtherUniqueId || testContext.parent.officeAccountRequestEdipi) {
    return true;
  }

  return false;
};

// It is TRANSCOM policy that an individual person may only be either a TIO or TOO, never both.
const validateOnlyOneTransportationOfficerRole = (value, testContext) => {
  const { taskOrderingOfficerCheckBox, taskInvoicingOfficerCheckBox } = testContext.parent;
  if (taskOrderingOfficerCheckBox && taskInvoicingOfficerCheckBox) {
    return false;
  }
  return true;
};

export const officeAccountRequestSchema = Yup.object().shape({
  officeAccountRequestFirstName: Yup.string()
    .matches(/^[A-Za-z]+$/, noNumericAllowedErrorMsg)
    .required('Required'),
  officeAccountRequestMiddleInitial: Yup.string()
    .matches(/^[A-Z]$/, middleInitialErrorMsg)
    .optional(),
  officeAccountRequestLastName: Yup.string()
    .matches(/^[A-Za-z]+$/, noNumericAllowedErrorMsg)
    .required('Required'),
  officeAccountRequestEdipi: Yup.string()
    .min(10, edipiMaxErrorMsg)
    .max(10, edipiMaxErrorMsg)
    .matches(/^[0-9]*$/, numericOnlyErrorMsg)
    .test('officeAccountRequestEdipi', 'Required if not using other unique identifier', validateEdipi),
  officeAccountRequestOtherUniqueId: Yup.string()
    .matches(/^[A-Za-z0-9]+$/, otherUniqueIdErrorMsg)
    .test('officeAccountRequestOtherUniqueId', 'Required if not using DODID#', validateOtherUniqueId),
  officeAccountRequestTelephone: phoneSchema.required('Required'),
  officeAccountRequestEmail: OfficeAccountRequestEmailSchema.required('Required'),
  officeAccountTransportationOffice: Yup.object().required('Required'),
  taskOrderingOfficerCheckBox: Yup.bool()
    .test('roleRequestedRequired', 'You must select at least one role.', validateRoleRequestedMethod)
    .test(
      'onlyOneTransportationOfficerRole',
      'You cannot select both Task Ordering Officer and Task Invoicing Officer. This is a policy managed by USTRANSCOM.',
      validateOnlyOneTransportationOfficerRole,
    ),
  taskInvoicingOfficerCheckBox: Yup.bool()
    .test('roleRequestedRequired', 'You must select at least one role.', validateRoleRequestedMethod)
    .test(
      'onlyOneTransportationOfficerRole',
      'You cannot select both Task Ordering Officer and Task Invoicing Officer. This is a policy managed by USTRANSCOM.',
      validateOnlyOneTransportationOfficerRole,
    ),
  servicesCounselorCheckBox: Yup.bool().test(
    'roleRequestedRequired',
    'You must select at least one role.',
    validateRoleRequestedMethod,
  ),
  transportationContractingOfficerCheckBox: Yup.bool().test(
    'roleRequestedRequired',
    'You must select at least one role.',
    validateRoleRequestedMethod,
  ),
  qualityAssuranceEvaluatorCheckBox: Yup.bool().test(
    'roleRequestedRequired',
    'You must select at least one role.',
    validateRoleRequestedMethod,
  ),
  headquartersCheckBox: Yup.bool().test(
    'roleRequestedRequired',
    'You must select at least one role.',
    validateRoleRequestedMethod,
  ),
  customerSupportRepresentativeCheckBox: Yup.bool().test(
    'roleRequestedRequired',
    'You must select at least one role.',
    validateRoleRequestedMethod,
  ),
  governmentSurveillanceRepresentativeCheckbox: Yup.bool().test(
    'roleRequestedRequired',
    'You must select at least one role.',
    validateRoleRequestedMethod,
  ),
});

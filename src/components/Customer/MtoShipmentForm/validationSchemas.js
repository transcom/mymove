import * as Yup from 'yup';

export const AgentSchema = Yup.object().shape({
  firstName: Yup.string(),
  lastName: Yup.string(),
  phone: Yup.string().matches(/^[2-9]\d{2}\d{3}\d{4}$/, 'Must be valid phone number'),
  email: Yup.string().email('Must be valid email'),
});

export const RequiredAddressSchema = Yup.object().shape({
  street_address_1: Yup.string().required('Required'),
  street_address_2: Yup.string(),
  city: Yup.string().required('Required'),
  state: Yup.string().length(2, 'Must use state abbreviation').required('Required'),
  postal_code: Yup.string()
    // RA Summary: eslint - security/detect-unsafe-regex - Denial of Service: Regular Expression
    // RA: Untrusted data is passed to the application and used as a regular expression. This can cause the thread to overconsume CPU resources.
    // RA: Line used for validating a zip code
    // RA: The regex pattern is a constant string set at compile-time and it is bounded. Therefore, it is not a risk
    // RA Developer Status: False Positive
    // RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
    // RA Validator: jneuner@mitre.org
    // RA Modified Severity:
    // eslint-disable-next-line security/detect-unsafe-regex
    .matches(/^(\d{5}([-]\d{4})?)$/, 'Must be valid zip code')
    .required('Required'),
});

export const OptionalAddressSchema = Yup.object().shape({
  street_address_1: Yup.string(),
  street_address_2: Yup.string(),
  city: Yup.string(),
  state: Yup.string().length(2, 'Must use state abbreviation'),
  postal_code: Yup.string()
    // RA Summary: eslint - security/detect-unsafe-regex - Denial of Service: Regular Expression
    // RA: Untrusted data is passed to the application and used as a regular expression. This can cause the thread to overconsume CPU resources.
    // RA: Line used for validating a zip code
    // RA: The regex pattern is a constant string set at compile-time and it is bounded. Therefore, it is not a risk
    // RA Developer Status: False Positive
    // RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
    // RA Validator: jneuner@mitre.org
    // RA Modified Severity:
    // eslint-disable-next-line security/detect-unsafe-regex
    .matches(/^(\d{5}([-]\d{4})?)$/, 'Must be valid zip code'),
});

export const RequiredPlaceSchema = Yup.object().shape({
  address: RequiredAddressSchema,
  agent: AgentSchema,
});

export const OptionalPlaceSchema = Yup.object().shape({
  address: OptionalAddressSchema,
  agent: AgentSchema,
});

export default { AgentSchema, RequiredAddressSchema, OptionalAddressSchema, RequiredPlaceSchema, OptionalPlaceSchema };

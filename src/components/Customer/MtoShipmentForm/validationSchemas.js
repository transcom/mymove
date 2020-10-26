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
    //  security/detect-unsafe-regex
    .matches(/^(\d{5}([-]\d{4})?)$/, 'Must be valid zip code')
    .required('Required'),
});

export const OptionalAddressSchema = Yup.object().shape({
  street_address_1: Yup.string(),
  street_address_2: Yup.string(),
  city: Yup.string(),
  state: Yup.string().length(2, 'Must use state abbreviation'),
  postal_code: Yup.string()
    //  security/detect-unsafe-regex
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

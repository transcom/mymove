import * as Yup from 'yup';

import { requiredAddressSchema } from 'utils/validation';
import { OptionalAddressSchema } from 'components/Shared/MtoShipmentForm/validationSchemas';

export const AgentSchema = Yup.object().shape({
  firstName: Yup.string(),
  lastName: Yup.string(),
  phone: Yup.string().matches(/^[2-9]\d{2}-\d{3}-\d{4}$/, 'Must be valid phone number'),
  email: Yup.string().email('Must be valid email'),
});

export const RequiredPlaceSchema = Yup.object().shape({
  address: requiredAddressSchema,
  agent: AgentSchema,
});

export const OptionalPlaceSchema = Yup.object().shape({
  address: OptionalAddressSchema,
  agent: AgentSchema,
});

export const AdditionalAddressSchema = Yup.object().shape({
  address: OptionalAddressSchema,
});

export const StorageFacilityAddressSchema = Yup.object().shape({
  address: requiredAddressSchema,
  lotNumber: Yup.string(),
  facilityName: Yup.string().required('Required'),
  phone: Yup.string().matches(/^[2-9]\d{2}-\d{3}-\d{4}$/, 'Must be valid phone number'),
  email: Yup.string().email('Must be valid email'),
});

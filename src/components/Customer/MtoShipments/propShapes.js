import { string, shape } from 'prop-types';

export const simpleAddressShape = shape({
  city: string,
  state: string,
  postal_code: string,
});

export const fullAddressShape = shape({
  ...simpleAddressShape,
  street_address_1: string,
  street_address_2: string,
});

export const agentShape = shape({
  firstName: string,
  lastName: string,
  phone: string,
  email: string,
  agentType: string,
});

export default { simpleAddressShape, fullAddressShape, agentShape };

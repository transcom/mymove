import { faker } from '@faker-js/faker';
import { build, oneOf, perBuild } from '@jackfranklin/test-data-bot';

import { fake, getInternalSpec } from 'utils/test/factories/base';

export const ADDRESS_FIELDS = {
  ID: 'id',
  STREET_ADDRESS_1: 'streetAddress1',
  STREET_ADDRESS_2: 'streetAddress2',
  CITY: 'city',
  STATE: 'state',
  POSTAL_CODE: 'postalCode',
  COUNTRY: 'country',
  ETAG: 'eTag',
};

const addressFactory = build({
  fields: {
    streetAddress1: fake((f) => f.address.streetAddress()),
    streetAddress2: fake((f) => f.address.secondaryAddress()),
    // left out streetAddress3 since we don't even let users input that line...
    city: fake((f) => f.address.city()),
    [ADDRESS_FIELDS.STATE]: perBuild(() => {
      const spec = getInternalSpec();
      return oneOf(...spec.definitions.Address.properties.state.enum).call();
    }),
    country: 'US', // Likely change once we support more than just OCONUS moves.
  },
  postBuild: (address) => {
    address.postalCode = faker.address.zipCodeByState(address.state);

    return address;
  },
});

export default addressFactory;

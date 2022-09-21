import { faker } from '@faker-js/faker';
import { build } from '@jackfranklin/test-data-bot';

import { idHelper, stateHelper } from './helpers';

import { fake } from 'utils/test/factories/base';

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
    [ADDRESS_FIELDS.ID]: fake(idHelper),
    [ADDRESS_FIELDS.STREET_ADDRESS_1]: fake((f) => f.address.streetAddress()),
    [ADDRESS_FIELDS.STREET_ADDRESS_2]: fake((f) => f.address.secondaryAddress()),
    // left out streetAddress3 since we don't even let users input that line...
    [ADDRESS_FIELDS.CITY]: fake((f) => f.address.city()),

    [ADDRESS_FIELDS.STATE]: fake(stateHelper),
    [ADDRESS_FIELDS.COUNTRY]: 'US', // Likely change once we support more than just OCONUS moves.
  },
  postBuild: (address) => {
    address.postalCode = faker.address.zipCodeByState(address.state);

    return address;
  },
});

export default addressFactory;

import { faker } from '@faker-js/faker';

import { idHelper, stateHelper } from './helpers';

import { baseFactory, BASE_FIELDS, fake, getInternalSpec } from 'utils/test/factories/base';

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

export const ADDRESS_TRAITS = {
  ONLY_REQUIRED_FIELDS: 'onlyRequiredFields',
};

const addressFactory = (params) => {
  return baseFactory({
    [BASE_FIELDS.FIELDS]: {
      [ADDRESS_FIELDS.ID]: fake(idHelper),
      [ADDRESS_FIELDS.STREET_ADDRESS_1]: fake((f) => f.address.streetAddress()),
      [ADDRESS_FIELDS.STREET_ADDRESS_2]: fake((f) => f.address.secondaryAddress()),
      // left out streetAddress3 since we don't even let users input that line...
      [ADDRESS_FIELDS.CITY]: fake((f) => f.address.city()),
      [ADDRESS_FIELDS.STATE]: fake(stateHelper),
      [ADDRESS_FIELDS.COUNTRY]: 'US', // Likely change once we support more than just CONUS moves: JIRA ticket MB-13996
    },
    postBuild: (address) => {
      address.postalCode = faker.address.zipCodeByState(address.state);

      return address;
    },
    [BASE_FIELDS.TRAITS]: {
      [ADDRESS_TRAITS.ONLY_REQUIRED_FIELDS]: {
        postBuild: (address) => {
          const spec = getInternalSpec();

          const requiredFields = new Set(spec.definitions.Address.required);

          Object.values(ADDRESS_FIELDS).forEach((field) => {
            if (!requiredFields.has(field)) {
              delete address[field];
            }
          });

          return address;
        },
      },
    },
    ...params,
  });
};

export default addressFactory;

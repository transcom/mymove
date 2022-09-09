import { faker } from '@faker-js/faker';
import { build, oneOf, perBuild } from '@jackfranklin/test-data-bot';
import { v4 } from 'uuid';

import { fake, getInternalSpec } from 'utils/test/factories/base';

export const ADDRESS_TRAITS = {
  ONLY_BASIC_ADDRESS: 'onlyBasicAddress',
  ONLY_REQUIRED_FIELDS: 'onlyRequiredFields',
};

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

export const addressFactory = build({
  fields: {
    [ADDRESS_FIELDS.ID]: perBuild(() => v4()),
    [ADDRESS_FIELDS.STREET_ADDRESS_1]: fake((f) => f.address.streetAddress()),
    [ADDRESS_FIELDS.STREET_ADDRESS_2]: fake((f) => f.address.secondaryAddress()),
    // left out streetAddress3 since we don't even let users input that line...
    [ADDRESS_FIELDS.CITY]: fake((f) => f.address.city()),
    [ADDRESS_FIELDS.STATE]: perBuild(() => {
      const spec = getInternalSpec();
      return oneOf(...spec.definitions.Address.properties.state.enum).call();
    }),
    [ADDRESS_FIELDS.POSTAL_CODE]: '',
    [ADDRESS_FIELDS.COUNTRY]: 'US', // Likely change once we support more than just CONUS moves.
    [ADDRESS_FIELDS.ETAG]: perBuild(() => window.btoa(new Date().toISOString())),
  },
  traits: {
    [ADDRESS_TRAITS.ONLY_BASIC_ADDRESS]: {
      postBuild: (address) => {
        const extraFields = [ADDRESS_FIELDS.ID, ADDRESS_FIELDS.COUNTRY, ADDRESS_FIELDS.ETAG];
        extraFields.forEach((field) => {
          delete address[field];
        });

        return address;
      },
    },
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
  postBuild: (address) => {
    if (address.postalCode === '' && address.state !== '') {
      address.postalCode = faker.address.zipCodeByState(address.state);

      // The `zipCodeByState` function uses ranges of numbers for each state to generate the zip, but since some state
      // zip codes start with 0's, they get stripped out. This adds the zero back to the front if needed.
      address.postalCode = address.postalCode.padStart(5, '0');
    }

    return address;
  },
});

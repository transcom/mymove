import { faker } from '@faker-js/faker';
import { build } from '@jackfranklin/test-data-bot';

import { fake } from 'utils/test/factories/base';

export const BLANK_ADDRESS = 'blank';
export const ADDRESS_WITHOUT_COUNTRY = 'omitCountry';

export const ADDRESS_FIELDS = {
  STREET_ADDRESS_1: 'streetAddress1',
  STREET_ADDRESS_2: 'streetAddress2',
  CITY: 'city',
  STATE: 'state',
  POSTAL_CODE: 'postalCode',
  COUNTRY: 'country',
};

export const addressFactory = build({
  fields: {
    [ADDRESS_FIELDS.STREET_ADDRESS_1]: fake((f) => f.address.streetAddress()),
    [ADDRESS_FIELDS.STREET_ADDRESS_2]: fake((f) => f.address.secondaryAddress()),
    // left out streetAddress3 since we don't even let users input that line...
    [ADDRESS_FIELDS.CITY]: fake((f) => f.address.city()),
    [ADDRESS_FIELDS.STATE]: fake((f) => f.address.stateAbbr()),
    [ADDRESS_FIELDS.POSTAL_CODE]: '',
    [ADDRESS_FIELDS.COUNTRY]: 'US', // Likely change once we support more than just CONUS moves.
  },
  traits: {
    [BLANK_ADDRESS]: {
      overrides: {
        [ADDRESS_FIELDS.STREET_ADDRESS_1]: '',
        [ADDRESS_FIELDS.STREET_ADDRESS_2]: '',
        [ADDRESS_FIELDS.CITY]: '',
        [ADDRESS_FIELDS.STATE]: '',
        [ADDRESS_FIELDS.POSTAL_CODE]: '',
        [ADDRESS_FIELDS.COUNTRY]: '',
      },
    },
    [ADDRESS_WITHOUT_COUNTRY]: {
      postBuild: (address) => {
        delete address.country;

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

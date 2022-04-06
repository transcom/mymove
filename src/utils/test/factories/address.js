import { faker } from '@faker-js/faker';
import { build } from '@jackfranklin/test-data-bot';

import fake from 'utils/test/factories/base';

const addressBuilder = build({
  fields: {
    streetAddress1: fake((f) => f.address.streetAddress()),
    streetAddress2: fake((f) => f.address.secondaryAddress()),
    // left out streetAddress3 since we don't even let users input that line...
    city: fake((f) => f.address.city()),
    state: fake((f) => f.address.stateAbbr()),
    country: 'US', // Likely change once we support more than just OCONUS moves.
  },
  postBuild: (address) => {
    address.postalCode = faker.address.zipCodeByState(address.state);

    return address;
  },
});

export default addressBuilder;

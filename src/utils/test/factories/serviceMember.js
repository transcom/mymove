import { faker } from '@faker-js/faker';
import { bool, build, oneOf, perBuild } from '@jackfranklin/test-data-bot';
import { v4 as uuidv4 } from 'uuid';

import serviceMemberAgencies from 'content/serviceMemberAgencies';
import { ORDERS_RANK_OPTIONS } from 'constants/orders';
import WEIGHT_ENTITLEMENTS from 'constants/weightEntitlements';
import fake from 'utils/test/factories/base';
import addressBuilder from 'utils/test/factories/address';

export const PHONE_FORMAT = '###-###-####';

export const serviceMemberBuilder = build({
  fields: {
    id: perBuild(uuidv4),
    affiliation: oneOf(...Object.keys(serviceMemberAgencies)),
    edipi: fake((f) => f.datatype.number({ min: 1000000000, max: 9999999999 }).toString()),
    rank: oneOf(...Object.keys(ORDERS_RANK_OPTIONS)),
    first_name: fake((f) => f.name.firstName()),
    middle_name: fake((f) => f.name.middleName()),
    last_name: fake((f) => f.name.lastName()),
    telephone: fake((f) => f.phone.phoneNumber(PHONE_FORMAT)),
    secondary_telephone: fake((f) => f.phone.phoneNumber(PHONE_FORMAT)),
    email_is_preferred: bool(),
    phone_is_preferred: bool(),
    residential_address: perBuild(addressBuilder),
  },
  postBuild: (serviceMember) => {
    // These packages don't seem to have a way to do a weighed chance, so here's a way of doing having 1 in 4 chance of having a suffix
    const setSuffix = faker.random.arrayElement([true, false, false, false]);

    if (setSuffix) {
      serviceMember.suffix = faker.name.suffix();
    }

    serviceMember.personal_email = faker.internet.exampleEmail(serviceMember.first_name, serviceMember.last_name);

    // Need at least one to be preferred...
    if (!serviceMember.email_is_preferred && !serviceMember.phone_is_preferred) {
      const preferredContactMethod = faker.random.arrayElement(['email_is_preferred', 'phone_is_preferred']);

      serviceMember[preferredContactMethod] = true;
    }

    serviceMember.backup_mailing_address = addressBuilder({
      overrides: {
        city: serviceMember.residential_address.city,
        state: serviceMember.residential_address.state,
      },
    });

    serviceMember.weight_allotment = WEIGHT_ENTITLEMENTS[serviceMember.rank];

    return serviceMember;
  },
});

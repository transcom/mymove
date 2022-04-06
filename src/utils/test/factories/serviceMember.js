import { build, perBuild } from '@jackfranklin/test-data-bot';
import { v4 as uuidv4 } from 'uuid';

import fake from 'utils/test/factories/base';

const serviceMemberBuilder = build('ServiceMember', {
  fields: {
    id: perBuild(uuidv4),
    edipi: fake((faker) => faker.datatype.number({ min: 1000000000, max: 9999999999 })),
    first_name: fake((faker) => faker.name.firstName()),
    middle_name: fake((faker) => faker.name.middleName()),
    last_name: fake((faker) => faker.name.lastName()),
  },
});

export default serviceMemberBuilder;

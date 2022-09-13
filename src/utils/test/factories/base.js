import { faker } from '@faker-js/faker';
import { perBuild } from '@jackfranklin/test-data-bot';

faker.setLocale('en_US');

const fake = (callback) => {
  return perBuild(() => callback(faker));
};

export default fake;

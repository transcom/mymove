import { oneOf } from '@jackfranklin/test-data-bot';
import { faker } from '@faker-js/faker';
import { v4 as uuidv4 } from 'uuid';

import { getInternalSpec } from './base';

import serviceMemberAgencies from 'content/serviceMemberAgencies';

const PHONE_FORMAT = '###-###-####';
const spec = getInternalSpec();

export const agencyHelper = () => oneOf(...Object.keys(serviceMemberAgencies)).call();
export const createdAtHelper = (f) => f.date.recent();
export const gblocHelper = (f) => f.random.alpha({ count: 4, casing: 'upper' });
export const idHelper = () => uuidv4();
export const phoneHelper = (f) => f.phone.number(PHONE_FORMAT);
export const placeNameHelper = (f) => f.lorem.words(Math.random() * 4 + 1);
export const dutyLocationNameHelper = (f) => `${f.lorem.words(Math.random() * 4 + 1)} AFB`;
export const stateHelper = () => oneOf(...spec.definitions.Address.properties.state.enum).call();

export const updatedAtFromCreatedAt = (createdAt) => faker.date.between(createdAt, Date.now());

export default {
  agencyHelper,
  createdAtHelper,
  gblocHelper,
  idHelper,
  phoneHelper,
  placeNameHelper,
  dutyLocationNameHelper,
  stateHelper,
};

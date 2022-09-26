import addressFactory, { ADDRESS_FIELDS } from './address';
import { baseFactory, BASE_FIELDS, fake } from './base';
import transportationOfficeFactory, { TRANSPORTATION_OFFICE_FIELDS } from './transportationOffice';
import * as helpers from './helpers';

export const DUTY_LOCATION_FIELDS = {
  ADDRESS: 'address',
  ADDRESS_ID: 'addressId',
  AFFILIATION: 'affiliation',
  ID: 'id',
  NAME: 'name',
  TRANSPORTATION_OFFICE: 'transportationOffice',
  TRANSPORTATION_OFFICE_ID: 'transportationOfficeId',
  CREATED_AT: 'createdAt',
  UPDATED_AT: 'updatedAt',
};

const dutyLocationFactory = (params) => {
  return baseFactory({
    [BASE_FIELDS.FIELDS]: {
      [DUTY_LOCATION_FIELDS.ID]: fake(helpers.idHelper),
      [DUTY_LOCATION_FIELDS.ADDRESS]: (addressParams) => addressFactory(addressParams),
      [DUTY_LOCATION_FIELDS.AFFILIATION]: fake(helpers.agencyHelper),
      [DUTY_LOCATION_FIELDS.NAME]: fake((f) => `${f.lorem.words(Math.random() * 4 + 1)} AFB`),
      [DUTY_LOCATION_FIELDS.TRANSPORTATION_OFFICE]: (toParams) => transportationOfficeFactory(toParams),
      [DUTY_LOCATION_FIELDS.CREATED_AT]: fake(helpers.createdAtHelper),
    },
    [BASE_FIELDS.POST_BUILD]: (dutyLocation) => {
      dutyLocation[DUTY_LOCATION_FIELDS.ADDRESS][ADDRESS_FIELDS.ID] = dutyLocation.address.id;
      dutyLocation[DUTY_LOCATION_FIELDS.TRANSPORTATION_OFFICE][TRANSPORTATION_OFFICE_FIELDS.ID] =
        dutyLocation.transportationOffice.id;
      dutyLocation[DUTY_LOCATION_FIELDS.UPDATED_AT] = helpers.updatedAtFromCreatedAt(dutyLocation.createdAt);
    },
    ...params,
  });
};

export default dutyLocationFactory;

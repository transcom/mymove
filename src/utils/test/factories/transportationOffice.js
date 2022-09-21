import addressFactory from './address';
import { baseFactory, BASE_FIELDS, fake } from './base';
import { gblocHelper, idHelper, phoneHelper, placeNameHelper } from './helpers';

export const TRANSPORTATION_OFFICE_FIELDS = {
  ID: 'id',
  NAME: 'name',
  ADDRESS: 'address',
  PHONE_LINES: 'phone_lines',
  GBLOC: 'gbloc',
  LATITUDE: 'latitude',
  LONGITUDE: 'longitude',
};

const transportationOfficeFactory = (params) => {
  const getPhoneLines = () => {
    const numberOfLines = Math.random(4);
    const lines = [];
    for (let i = 0; i < numberOfLines; i += 1) {
      lines.push(fake(phoneHelper));
    }
    return lines;
  };
  return baseFactory({
    [BASE_FIELDS.FIELDS]: {
      [TRANSPORTATION_OFFICE_FIELDS.ID]: fake(idHelper),
      [TRANSPORTATION_OFFICE_FIELDS.NAME]: fake(placeNameHelper),
      [TRANSPORTATION_OFFICE_FIELDS.ADDRESS]: fake((addressParams) => addressFactory(addressParams)),
      [TRANSPORTATION_OFFICE_FIELDS.PHONE_LINES]: getPhoneLines,
      [TRANSPORTATION_OFFICE_FIELDS.GBLOC]: fake(gblocHelper),
      [TRANSPORTATION_OFFICE_FIELDS.LATITUDE]: fake((f) => f.address.latitude()),
      [TRANSPORTATION_OFFICE_FIELDS.LONGITUDE]: fake((f) => f.address.longitude()),
    },
    ...params,
  });
};

export default transportationOfficeFactory;

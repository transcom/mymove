import addressFactory from './address';
import { baseFactory } from './base';

export const DUTY_LOCATION_FIELDS = {
  ADDRESS: 'address',
};

const dutyLocationFactory = (params) => {
  return baseFactory({
    fields: {
      address: (subparams) => addressFactory(subparams),
    },
    ...params,
  });
};

export default dutyLocationFactory;

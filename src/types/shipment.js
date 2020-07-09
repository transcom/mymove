import PropTypes from 'prop-types';

import { SHIPMENT_TYPE } from 'shared/constants';

// eslint-disable-next-line import/prefer-default-export
export const ShipmentTypeOneOf = PropTypes.oneOf([
  SHIPMENT_TYPE.HHG,
  SHIPMENT_TYPE.HHG_SHORTHAUL_DOMESTIC,
  SHIPMENT_TYPE.HHG_LONGHAUL_DOMESTIC,
  SHIPMENT_TYPE.NTS,
]);

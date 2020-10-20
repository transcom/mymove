import PropTypes from 'prop-types';

import { SHIPMENT_OPTIONS, MOVE_TYPES } from 'shared/constants';

// eslint-disable-next-line import/prefer-default-export
export const ShipmentOptionsOneOf = PropTypes.oneOf([
  SHIPMENT_OPTIONS.HHG,
  SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC,
  SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
  SHIPMENT_OPTIONS.NTS,
  MOVE_TYPES.NTS,
  SHIPMENT_OPTIONS.NTSR,
  MOVE_TYPES.PPM,
]);

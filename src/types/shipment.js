import PropTypes from 'prop-types';

import { SHIPMENT_OPTIONS } from 'shared/constants';

// eslint-disable-next-line import/prefer-default-export
export const ShipmentOptionsOneOf = PropTypes.oneOf([
  SHIPMENT_OPTIONS.HHG,
  SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC,
  SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
  SHIPMENT_OPTIONS.NTS,
  SHIPMENT_OPTIONS.NTSR,
  SHIPMENT_OPTIONS.PPM,
]);

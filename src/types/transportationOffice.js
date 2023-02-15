/* eslint-disable import/prefer-default-export */
import PropTypes from 'prop-types';

import { AddressShape } from './address';

export const TransportationOfficeShape = PropTypes.shape({
  address: AddressShape,
  address_id: PropTypes.string,
  gbloc: PropTypes.string,
  created_at: PropTypes.string,
  id: PropTypes.string,
  name: PropTypes.string,
  updated_at: PropTypes.string,
});

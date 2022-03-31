/* eslint-disable import/prefer-default-export */
import PropTypes from 'prop-types';

import { AddressShape } from './address';

export const DutyLocationShape = PropTypes.shape({
  address: AddressShape,
  address_id: PropTypes.string,
  affiliation: PropTypes.string,
  created_at: PropTypes.string,
  id: PropTypes.string,
  name: PropTypes.string,
  updated_at: PropTypes.string,
});

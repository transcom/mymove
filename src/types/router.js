/* eslint-disable import/prefer-default-export */

import PropTypes from 'prop-types';

export const MatchShape = PropTypes.shape({
  params: PropTypes.object,
  isExact: PropTypes.bool,
  path: PropTypes.string,
  url: PropTypes.string,
});

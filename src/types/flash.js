/* eslint-disable import/prefer-default-export */
import PropTypes from 'prop-types';

export const FlashMessageShape = PropTypes.shape({
  type: PropTypes.string.isRequired,
  title: PropTypes.string,
  message: PropTypes.string.isRequired,
  key: PropTypes.string.isRequired,
});

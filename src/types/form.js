/* eslint-disable import/prefer-default-export */
import PropTypes from 'prop-types';

export const DropdownArrayOf = PropTypes.arrayOf(
  PropTypes.shape({
    key: PropTypes.string.isRequired,
    value: PropTypes.string.isRequired,
  }),
);

import PropTypes from 'prop-types';

export const AddressShape = PropTypes.shape({
  street_address_1: PropTypes.string,
  street_address_2: PropTypes.string,
  street_address_3: PropTypes.string,
  city: PropTypes.string,
  state: PropTypes.string,
  postal_code: PropTypes.string,
  country: PropTypes.string,
});

import React from 'react';
import PropTypes from 'prop-types';

export default function Address({ address }) {
  if (!address) {
    return null;
  }

  return (
    <div className="address">
      <div>{address.street_address_1}</div>
      {address.street_address_2 && <div>{address.street_address_2}</div>}
      <div>
        {address.city}, {address.state} {address.postal_code}
      </div>
    </div>
  );
}

Address.propTypes = {
  address: PropTypes.object,
};

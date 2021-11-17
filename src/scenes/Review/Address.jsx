import React from 'react';
import PropTypes from 'prop-types';

export default function Address({ address }) {
  if (!address) {
    return null;
  }

  return (
    <div className="address">
      <div>{address.streetAddress1}</div>
      {address.streetAddress2 && <div>{address.streetAddress2}</div>}
      <div>
        {address.city}, {address.state} {address.postalCode}
      </div>
    </div>
  );
}

Address.propTypes = {
  address: PropTypes.object,
};

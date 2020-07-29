/* eslint-disable camelcase */
import React from 'react';

export default function formatAddress(address) {
  const { street_address_1, city, state, postal_code } = address;
  return (
    <>
      {street_address_1 && (
        <>
          {street_address_1}
          <br />
        </>
      )}
      {city ? `${city}, ${state} ${postal_code}` : postal_code}
    </>
  );
}

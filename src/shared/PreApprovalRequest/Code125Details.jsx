import React, { Fragment } from 'react';
import { formatDate } from 'shared/formatters';

export const Code125Details = props => {
  const row = props.shipmentLineItem;
  return (
    <td data-cy={`${row.tariff400ng_item.code}-details`}>
      {row.reason} <br />
      {`Date of service: ${formatDate(row.date)}`} <br />
      {row.time && (
        <Fragment>
          {`Time of service: ${row.time}`} <br />
        </Fragment>
      )}
      {row.address.street_address_1} <br />
      {row.address.street_address_2 && (
        <Fragment>
          {row.address.street_address_2} <br />
        </Fragment>
      )}
      {row.address.city}, {row.address.state} {row.address.postal_code} <br />
    </td>
  );
};

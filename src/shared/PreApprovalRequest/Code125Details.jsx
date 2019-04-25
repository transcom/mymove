import React, { Fragment } from 'react';
import { formatDate } from 'shared/formatters';
import { AddressDisplay } from 'shared/Address';

export const Code125Details = props => {
  const row = props.shipmentLineItem;
  return (
    <td details-cy={`${row.tariff400ng_item.code}-details`}>
      {row.reason} <br />
      {`Date of service: ${formatDate(row.date)}`} <br />
      {row.time && (
        <Fragment>
          {`Time of service: ${row.time}`} <br />
        </Fragment>
      )}
      <AddressDisplay address={row.address} />
    </td>
  );
};

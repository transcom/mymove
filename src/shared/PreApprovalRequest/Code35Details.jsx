import React from 'react';

export const Code35Details = props => {
  const row = props.shipmentLineItem;
  return (
    <td details-cy={`${row.tariff400ng_item.code}-details`}>
      {row.description} <br />
      {row.reason} <br />
      Est. not to exceed: ${row.estimate_amount_cents} <br />
      Actual Amount: {row.estimate_amount_cents || `--`}
    </td>
  );
};

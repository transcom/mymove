import React from 'react';
import { formatCents } from 'shared/formatters';

export const Code226Details = props => {
  const row = props.shipmentLineItem;
  const actAmtValue = formatCents(row.actual_amount_cents);
  return (
    <td details-cy={`${row.tariff400ng_item.code}-details`}>
      {row.description} <br />
      {row.reason} <br />
      {`$${actAmtValue}`} <br />
      {row.notes}
    </td>
  );
};

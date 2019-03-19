import React from 'react';
import { formatCents } from 'shared/formatters';

export const Code35Details = props => {
  const row = props.shipmentLineItem;
  return (
    <td details-cy={`${row.tariff400ng_item.code}-details`}>
      {row.description} <br />
      {row.reason} <br />
      Est. not to exceed: ${formatCents(row.estimate_amount_cents)} <br />
      Actual Amount: {`${formatCents(row.actual_amount_cents)}` || `--`} <br />
      {row.notes}
    </td>
  );
};

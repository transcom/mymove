import React from 'react';
import { formatFromBaseQuantity } from 'shared/formatters';

export const DefaultDetails = props => {
  const row = props.shipmentLineItem;
  return (
    <td details-cy={`${row.tariff400ng_item.code}-default-details`}>
      {formatFromBaseQuantity(row.quantity_1)} <br />
      {row.notes}
    </td>
  );
};

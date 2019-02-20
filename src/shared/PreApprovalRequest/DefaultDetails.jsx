import React from 'react';
import { formatFromBaseQuantity } from 'shared/formatters';

export const DefaultDetails = props => {
  const row = props.shipmentLineItem;
  return (
    <td align="left">
      {formatFromBaseQuantity(row.quantity_1)} <br />
      {row.notes}
    </td>
  );
};

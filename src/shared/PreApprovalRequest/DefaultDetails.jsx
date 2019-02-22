import React from 'react';
import { formatFromBaseQuantity } from 'shared/formatters';

export const DefaultDetails = props => {
  const row = props.shipmentLineItem;
  return (
    <td>
      {formatFromBaseQuantity(row.quantity_1)} <br />
      {row.notes}
    </td>
  );
};

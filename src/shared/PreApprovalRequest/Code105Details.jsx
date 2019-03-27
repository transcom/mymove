import React from 'react';
import { convertFromThousandthInchToInch } from 'shared/formatters';
import { displayBaseQuantityUnits } from 'shared/lineItems';

export const Code105Details = props => {
  const row = props.shipmentLineItem;
  let crateLengthinInches = convertFromThousandthInchToInch(row.crate_dimensions.length);
  let crateWidthinInches = convertFromThousandthInchToInch(row.crate_dimensions.width);
  let crateHeightinInches = convertFromThousandthInchToInch(row.crate_dimensions.height);
  let itemLengthinInches = convertFromThousandthInchToInch(row.item_dimensions.length);
  let itemWidthinInches = convertFromThousandthInchToInch(row.item_dimensions.width);
  let itemHeightinInches = convertFromThousandthInchToInch(row.item_dimensions.height);
  let crateCubicFeet = displayBaseQuantityUnits(row);

  let crateDetails = `Crate: ${crateLengthinInches}" x ${crateWidthinInches}" x ${crateHeightinInches}" (${crateCubicFeet})`;
  let ItemDetails = `Item: ${itemLengthinInches}" x ${itemWidthinInches}" x ${itemHeightinInches}"`;
  return (
    <td details-cy={`${row.tariff400ng_item.code}-details`}>
      {row.description} <br />
      {crateDetails} <br />
      {ItemDetails} <br />
      {row.notes}
    </td>
  );
};

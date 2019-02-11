import { truncateNumber, addCommasToNumberString, formatFromBaseQuantity } from 'shared/formatters';

export const displayBaseQuantityUnits = (item, scale) => {
  if (!item) return;

  const itemCode = item.tariff400ng_item.code;
  const itemQuantity = item.quantity_1;

  const lbsItems = ['105A', '105C', '135A', '135B'];

  if (lbsItems.includes(itemCode)) {
    const decimalPlaces = 0;
    const baseQuantity = truncateNumber(itemQuantity, decimalPlaces);
    return `${addCommasToNumberString(baseQuantity, decimalPlaces)} lbs`;
  }
  return formatFromBaseQuantity(itemQuantity);
};

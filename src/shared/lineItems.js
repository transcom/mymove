import {
  truncateNumber,
  addCommasToNumberString,
  formatFromBaseQuantity,
  convertFromBaseQuantity,
} from 'shared/formatters';
import { isNewAccessorial } from 'shared/preApprovals';

export const displayBaseQuantityUnits = (item, scale) => {
  if (!item) return;

  const itemCode = item.tariff400ng_item.code;
  const itemQuantity = item.quantity_1;
  const lbsItems = ['105A', '105C', '135A', '135B'];
  const cuFtItems = ['105B', '105E'];

  switch (true) {
    case lbsItems.includes(itemCode): {
      const decimalPlaces = 0;
      const convertedItemQuantity = convertFromBaseQuantity(itemQuantity);
      const baseQuantity = truncateNumber(convertedItemQuantity, decimalPlaces);
      return `${addCommasToNumberString(baseQuantity, decimalPlaces)} lbs`;
    }
    case cuFtItems.includes(itemCode) && isNewAccessorial(item): {
      const decimalPlaces = 2;
      const convertedItemQuantity = convertFromBaseQuantity(itemQuantity);
      const baseQuantity = truncateNumber(convertedItemQuantity, decimalPlaces);
      return `${addCommasToNumberString(baseQuantity, decimalPlaces)} cu ft`;
    }
    default:
      return formatFromBaseQuantity(itemQuantity);
  }
};

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
  const itemQuantity1 = item.quantity_1;
  const itemQuantity2 = item.quantity_2;
  const lbsItems = ['105A', '105C', '135A', '135B'];
  const cuFtItems = ['105B', '105E'];
  const lbsMiItems = ['LHS', '16A'];

  if (lbsItems.includes(itemCode)) {
    const decimalPlaces = 0;
    const convertedItemQuantity1 = convertFromBaseQuantity(itemQuantity1);
    const weight = truncateNumber(convertedItemQuantity1, decimalPlaces);
    return `${addCommasToNumberString(weight, decimalPlaces)} lbs`;
  } else if (lbsMiItems.includes(itemCode)) {
    const decimalPlaces = 0;
    const convertedItemQuantity1 = convertFromBaseQuantity(itemQuantity1);
    const convertedItemQuantity2 = convertFromBaseQuantity(itemQuantity2);
    const weight = truncateNumber(convertedItemQuantity1, decimalPlaces);
    const milage = truncateNumber(convertedItemQuantity2, decimalPlaces);
    return `${addCommasToNumberString(weight, decimalPlaces)} lbs, ${addCommasToNumberString(
      milage,
      decimalPlaces,
    )} mi`;
  } else if (cuFtItems.includes(itemCode) && isNewAccessorial(item)) {
    const decimalPlaces = 2;
    const convertedItemQuantity1 = convertFromBaseQuantity(itemQuantity1);
    const volume = truncateNumber(convertedItemQuantity1, decimalPlaces);
    return `${addCommasToNumberString(volume, decimalPlaces)} cu ft`;
  }
  return formatFromBaseQuantity(itemQuantity1);
};

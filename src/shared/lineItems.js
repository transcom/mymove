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

  if (isWeight(itemCode)) {
    const decimalPlaces = 0;
    const weight = convertTruncateAddCommas(itemQuantity1, decimalPlaces);
    return `${weight} lbs`;
  } else if (isWeightDistance(itemCode)) {
    const decimalPlaces = 0;
    const weight = convertTruncateAddCommas(itemQuantity1, decimalPlaces);
    const milage = convertTruncateAddCommas(itemQuantity2, decimalPlaces);
    return `${weight} lbs, ${milage} mi`;
  } else if (isVolume(itemCode) && isNewAccessorial(item)) {
    const decimalPlaces = 2;
    const volume = convertTruncateAddCommas(itemQuantity1, decimalPlaces);
    return `${volume} cu ft`;
  }
  return formatFromBaseQuantity(itemQuantity1);
};

function isWeight(itemCode) {
  const lbsItems = ['105A', '105C', '135A', '135B'];
  return lbsItems.includes(itemCode);
}

function isVolume(itemCode) {
  const cuFtItems = ['105B', '105E'];
  return cuFtItems.includes(itemCode);
}

function isWeightDistance(itemCode) {
  const lbsMiItems = ['LHS', '16A'];
  return lbsMiItems.includes(itemCode);
}

function convertTruncateAddCommas(value, decimalPlaces) {
  const convertedValue = convertFromBaseQuantity(value);
  const formattedValue = truncateNumber(convertedValue, decimalPlaces);
  return addCommasToNumberString(formattedValue, decimalPlaces);
}

import { SERVICE_ITEM_CODES } from 'constants/serviceItems';
import { serviceItemCodes } from 'content/serviceItems';

export const createServiceItemModelTypes = {
  MTOServiceItemOriginSIT: 'MTOServiceItemOriginSIT',
  MTOServiceItemDestSIT: 'MTOServiceItemDestSIT',
  MTOServiceItemShuttle: 'MTOServiceItemShuttle',
  MTOServiceItemDomesticCrating: 'MTOServiceItemDomesticCrating',
  MTOServiceItemDomesticStandaloneCrating: 'MTOServiceItemDomesticStandaloneCrating',
};

export const shuttleServiceItemCodeOptions = [
  { value: serviceItemCodes.DOSHUT, key: SERVICE_ITEM_CODES.DOSHUT },
  { value: serviceItemCodes.DDSHUT, key: SERVICE_ITEM_CODES.DDSHUT },
];

export const domesticCratingServiceItemCodeOptions = [
  { value: serviceItemCodes.DCRT, key: SERVICE_ITEM_CODES.DCRT },
  { value: serviceItemCodes.DUCRT, key: SERVICE_ITEM_CODES.DUCRT },
];

export const domesticStandaloneCratingServiceItemCodeOptions = [
  { value: serviceItemCodes.DCRTSA, key: SERVICE_ITEM_CODES.DCRTSA },
];

export default createServiceItemModelTypes;

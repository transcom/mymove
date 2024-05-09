import { SERVICE_ITEM_CODES } from 'constants/serviceItems';
import { serviceItemCodes } from 'content/serviceItems';

export const createServiceItemModelTypes = {
  MTOServiceItemOriginSIT: 'MTOServiceItemOriginSIT',
  MTOServiceItemDestSIT: 'MTOServiceItemDestSIT',
  MTOServiceItemShuttle: 'MTOServiceItemShuttle',
  MTOServiceItemDomesticCrating: 'MTOServiceItemDomesticCrating',
  MTOServiceItemStandaloneCrating: 'MTOServiceItemStandaloneCrating',
};

export const shuttleServiceItemCodeOptions = [
  { value: serviceItemCodes.DOSHUT, key: SERVICE_ITEM_CODES.DOSHUT },
  { value: serviceItemCodes.DDSHUT, key: SERVICE_ITEM_CODES.DDSHUT },
];

export const domesticCratingServiceItemCodeOptions = [
  { value: serviceItemCodes.DCRT, key: SERVICE_ITEM_CODES.DCRT },
  { value: serviceItemCodes.DUCRT, key: SERVICE_ITEM_CODES.DUCRT },
];

export const standaloneCratingServiceItemCodeOptions = [
  { value: serviceItemCodes.SCRT, key: SERVICE_ITEM_CODES.SCRT },
  { value: serviceItemCodes.SUCRT, key: SERVICE_ITEM_CODES.SUCRT },
];

export default createServiceItemModelTypes;

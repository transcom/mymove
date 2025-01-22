import { SERVICE_ITEM_CODES } from 'constants/serviceItems';
import { serviceItemCodes } from 'content/serviceItems';

export const createServiceItemModelTypes = {
  MTOServiceItemOriginSIT: 'MTOServiceItemOriginSIT',
  MTOServiceItemDestSIT: 'MTOServiceItemDestSIT',
  MTOServiceItemInternationalOriginSIT: 'MTOServiceItemInternationalOriginSIT',
  MTOServiceItemInternationalDestSIT: 'MTOServiceItemInternationalDestSIT',
  MTOServiceItemShuttle: 'MTOServiceItemShuttle',
  MTOServiceItemDomesticCrating: 'MTOServiceItemDomesticCrating',
  MTOServiceItemInternationalCrating: 'MTOServiceItemInternationalCrating',
};

export const shuttleServiceItemCodeOptions = [
  { value: serviceItemCodes.DOSHUT, key: SERVICE_ITEM_CODES.DOSHUT },
  { value: serviceItemCodes.DDSHUT, key: SERVICE_ITEM_CODES.DDSHUT },
];

export const domesticCratingServiceItemCodeOptions = [
  { value: serviceItemCodes.DCRT, key: SERVICE_ITEM_CODES.DCRT },
  { value: serviceItemCodes.DUCRT, key: SERVICE_ITEM_CODES.DUCRT },
];

export const internationalCratingServiceItemCodeOptions = [
  { value: serviceItemCodes.ICRT, key: SERVICE_ITEM_CODES.ICRT },
  { value: serviceItemCodes.IUCRT, key: SERVICE_ITEM_CODES.IUCRT },
];

export default createServiceItemModelTypes;

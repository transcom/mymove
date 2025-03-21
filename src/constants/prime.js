import { SERVICE_ITEM_CODES } from 'constants/serviceItems';
import { serviceItemCodes } from 'content/serviceItems';

export const createServiceItemModelTypes = {
  MTOServiceItemOriginSIT: 'MTOServiceItemOriginSIT',
  MTOServiceItemDestSIT: 'MTOServiceItemDestSIT',
  MTOServiceItemInternationalOriginSIT: 'MTOServiceItemInternationalOriginSIT',
  MTOServiceItemInternationalDestSIT: 'MTOServiceItemInternationalDestSIT',
  MTOServiceItemDomesticShuttle: 'MTOServiceItemDomesticShuttle',
  MTOServiceItemDomesticCrating: 'MTOServiceItemDomesticCrating',
  MTOServiceItemInternationalCrating: 'MTOServiceItemInternationalCrating',
  MTOServiceItemInternationalShuttle: 'MTOServiceItemInternationalShuttle',
};

export const domesticShuttleServiceItemCodeOptions = [
  { value: serviceItemCodes.DOSHUT, key: SERVICE_ITEM_CODES.DOSHUT },
  { value: serviceItemCodes.DDSHUT, key: SERVICE_ITEM_CODES.DDSHUT },
];

export const internationalShuttleServiceItemCodeOptions = [
  { value: serviceItemCodes.IOSHUT, key: SERVICE_ITEM_CODES.IOSHUT },
  { value: serviceItemCodes.IDSHUT, key: SERVICE_ITEM_CODES.IDSHUT },
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

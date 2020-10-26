import * as Yup from 'yup';

import { RequiredPlaceSchema, OptionalPlaceSchema } from './validationSchemas';

import { SHIPMENT_OPTIONS } from 'shared/constants';

const hhgShipmentSchema = Yup.object().shape({
  pickup: RequiredPlaceSchema,
  delivery: OptionalPlaceSchema,
  customerRemarks: Yup.string(),
});

const ntsShipmentSchema = Yup.object().shape({
  pickup: RequiredPlaceSchema,
  customerRemarks: Yup.string(),
});

const ntsReleaseShipmentSchema = Yup.object().shape({
  delivery: OptionalPlaceSchema,
  customerRemarks: Yup.string(),
});

export function getShipmentOptions(shipmentType) {
  switch (shipmentType) {
    case SHIPMENT_OPTIONS.HHG:
      return {
        schema: hhgShipmentSchema,
        showPickupFields: true,
        showDeliveryFields: true,
        displayName: 'HHG',
      };
    case SHIPMENT_OPTIONS.NTS:
      return {
        schema: ntsShipmentSchema,
        showPickupFields: true,
        showDeliveryFields: false,
        displayName: 'NTS',
      };
    case SHIPMENT_OPTIONS.NTSR:
      return {
        schema: ntsReleaseShipmentSchema,
        showPickupFields: false,
        showDeliveryFields: true,
        displayName: 'NTS-R',
      };
    default:
      throw new Error('unrecognized move type');
  }
}

export default { getShipmentOptions };

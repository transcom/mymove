import * as Yup from 'yup';

import { AdditionalAddressSchema, RequiredPlaceSchema, OptionalPlaceSchema } from './validationSchemas';

import { SHIPMENT_OPTIONS } from 'shared/constants';

const hhgShipmentSchema = Yup.object().shape({
  pickup: RequiredPlaceSchema,
  delivery: OptionalPlaceSchema,
  secondaryPickup: AdditionalAddressSchema,
  secondaryDelivery: AdditionalAddressSchema,
  customerRemarks: Yup.string(),
  counselorRemarks: Yup.string(),
});

const ntsShipmentSchema = Yup.object().shape({
  pickup: RequiredPlaceSchema,
  secondaryPickup: AdditionalAddressSchema,
  customerRemarks: Yup.string(),
});

const ntsReleaseShipmentSchema = Yup.object().shape({
  delivery: OptionalPlaceSchema,
  secondaryDelivery: AdditionalAddressSchema,
  customerRemarks: Yup.string(),
});

function getShipmentOptions(shipmentType) {
  switch (shipmentType) {
    case SHIPMENT_OPTIONS.HHG:
      return {
        schema: hhgShipmentSchema,
        showPickupFields: true,
        showDeliveryFields: true,
      };
    case SHIPMENT_OPTIONS.NTS:
      return {
        schema: ntsShipmentSchema,
        showPickupFields: true,
        showDeliveryFields: false,
      };
    case SHIPMENT_OPTIONS.NTSR:
      return {
        schema: ntsReleaseShipmentSchema,
        showPickupFields: false,
        showDeliveryFields: true,
      };
    default:
      throw new Error('unrecognized move type');
  }
}

export default getShipmentOptions;

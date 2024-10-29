import * as Yup from 'yup';

import {
  AdditionalAddressSchema,
  RequiredPlaceSchema,
  OptionalPlaceSchema,
  StorageFacilityAddressSchema,
} from './validationSchemas';

import { SHIPMENT_OPTIONS, SHIPMENT_TYPES } from 'shared/constants';
import { roleTypes } from 'constants/userRoles';

const hhgShipmentSchema = Yup.object().shape({
  pickup: RequiredPlaceSchema,
  delivery: OptionalPlaceSchema,
  secondaryPickup: AdditionalAddressSchema,
  secondaryDelivery: AdditionalAddressSchema,
  tertiaryPickup: AdditionalAddressSchema,
  tertiaryDelivery: AdditionalAddressSchema,
  customerRemarks: Yup.string(),
  counselorRemarks: Yup.string(),
});

const mobileHomeShipmentLocationSchema = Yup.object().shape({
  pickup: RequiredPlaceSchema,
  delivery: OptionalPlaceSchema,
  secondaryPickup: AdditionalAddressSchema,
  secondaryDelivery: AdditionalAddressSchema,
  tertiaryPickup: AdditionalAddressSchema,
  tertiaryDelivery: AdditionalAddressSchema,
  customerRemarks: Yup.string(),
  counselorRemarks: Yup.string(),
});

const boatShipmentLocationInfoSchema = Yup.object().shape({
  pickup: RequiredPlaceSchema,
  delivery: OptionalPlaceSchema,
  secondaryPickup: AdditionalAddressSchema,
  secondaryDelivery: AdditionalAddressSchema,
  tertiaryPickup: AdditionalAddressSchema,
  tertiaryDelivery: AdditionalAddressSchema,
  customerRemarks: Yup.string(),
  counselorRemarks: Yup.string(),
});

const ntsShipmentSchema = Yup.object().shape({
  pickup: RequiredPlaceSchema,
  secondaryPickup: AdditionalAddressSchema,
  tertiaryPickup: AdditionalAddressSchema,
  customerRemarks: Yup.string(),
  serviceOrderNumber: Yup.string().matches(/^[0-9a-zA-Z]+$/, 'Letters and numbers only'),
});

const ntsShipmentTOOSchema = Yup.object().shape({
  pickup: RequiredPlaceSchema,
  secondaryPickup: AdditionalAddressSchema,
  tertiaryPickup: AdditionalAddressSchema,
  serviceOrderNumber: Yup.string()
    .required('Required')
    .matches(/^[0-9a-zA-Z]+$/, 'Letters and numbers only'),
  storageFacility: StorageFacilityAddressSchema,
});

const ntsReleaseShipmentSchema = Yup.object().shape({
  delivery: RequiredPlaceSchema,
  secondaryDelivery: AdditionalAddressSchema,
  tertiaryDelivery: AdditionalAddressSchema,
  customerRemarks: Yup.string(),
});

const ntsReleaseShipmentCounselorSchema = Yup.object().shape({
  delivery: RequiredPlaceSchema,
  secondaryDelivery: AdditionalAddressSchema,
  tertiaryDelivery: AdditionalAddressSchema,
  counselorRemarks: Yup.string(),
  serviceOrderNumber: Yup.string().matches(/^[0-9a-zA-Z]+$/, 'Letters and numbers only'),
  storageFacility: StorageFacilityAddressSchema,
});

const ntsReleaseShipmentTOOSchema = Yup.object().shape({
  delivery: RequiredPlaceSchema,
  ntsRecordedWeight: Yup.string().required('Required'),
  secondaryDelivery: AdditionalAddressSchema,
  tertiaryDelivery: AdditionalAddressSchema,
  serviceOrderNumber: Yup.string()
    .required('Required')
    .matches(/^[0-9a-zA-Z]+$/, 'Letters and numbers only'),
  storageFacility: StorageFacilityAddressSchema,
});

function getShipmentOptions(shipmentType, userRole) {
  switch (shipmentType) {
    case SHIPMENT_OPTIONS.HHG:
      return {
        schema: hhgShipmentSchema,
        showPickupFields: true,
        showDeliveryFields: true,
      };

    case SHIPMENT_OPTIONS.MOBILE_HOME:
      return {
        schema: mobileHomeShipmentLocationSchema,
        showPickupFields: true,
        showDeliveryFields: true,
      };

    case SHIPMENT_OPTIONS.BOAT:
    case SHIPMENT_TYPES.BOAT_HAUL_AWAY:
    case SHIPMENT_TYPES.BOAT_TOW_AWAY:
      return {
        schema: boatShipmentLocationInfoSchema,
        showPickupFields: true,
        showDeliveryFields: true,
      };

    case SHIPMENT_OPTIONS.NTS:
      switch (userRole) {
        case roleTypes.TOO: {
          return {
            schema: ntsShipmentTOOSchema,
            showPickupFields: true,
            showDeliveryFields: true,
          };
        }

        default: {
          return {
            schema: ntsShipmentSchema,
            showPickupFields: true,
            showDeliveryFields: false,
          };
        }
      }

    case SHIPMENT_OPTIONS.NTSR:
      switch (userRole) {
        case roleTypes.CUSTOMER: {
          return {
            schema: ntsReleaseShipmentSchema,
            showPickupFields: false,
            showDeliveryFields: true,
          };
        }

        case roleTypes.SERVICES_COUNSELOR: {
          return {
            schema: ntsReleaseShipmentCounselorSchema,
            showPickupFields: false,
            showDeliveryFields: true,
          };
        }

        case roleTypes.TOO: {
          return {
            schema: ntsReleaseShipmentTOOSchema,
            showPickupFields: true,
            showDeliveryFields: true,
          };
        }

        default: {
          throw new Error('unrecognized user role type');
        }
      }

    case SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE:
      return {
        schema: hhgShipmentSchema,
        showPickupFields: true,
        showDeliveryFields: true,
      };

    default:
      throw new Error('unrecognized shipment type');
  }
}

export default getShipmentOptions;

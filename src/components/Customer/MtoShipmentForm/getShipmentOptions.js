import * as Yup from 'yup';

import {
  AdditionalAddressSchema,
  RequiredPlaceSchema,
  OptionalPlaceSchema,
  StorageFacilityAddressSchema,
} from './validationSchemas';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { roleTypes } from 'constants/userRoles';

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
  serviceOrderNumber: Yup.string().matches(/^[0-9a-zA-Z]+$/, 'Letters and numbers only'),
});

const ntsShipmentTOOSchema = Yup.object().shape({
  pickup: RequiredPlaceSchema,
  secondaryPickup: AdditionalAddressSchema,
  serviceOrderNumber: Yup.string()
    .required('Required')
    .matches(/^[0-9a-zA-Z]+$/, 'Letters and numbers only'),
  storageFacility: StorageFacilityAddressSchema,
});

const ntsReleaseShipmentSchema = Yup.object().shape({
  delivery: RequiredPlaceSchema,
  secondaryDelivery: AdditionalAddressSchema,
  customerRemarks: Yup.string(),
});

const ntsReleaseShipmentCounselorSchema = Yup.object().shape({
  delivery: RequiredPlaceSchema,
  secondaryDelivery: AdditionalAddressSchema,
  counselorRemarks: Yup.string(),
  serviceOrderNumber: Yup.string().matches(/^[0-9a-zA-Z]+$/, 'Letters and numbers only'),
  storageFacility: StorageFacilityAddressSchema,
});

const ntsReleaseShipmentTOOSchema = Yup.object().shape({
  delivery: RequiredPlaceSchema,
  ntsRecordedWeight: Yup.string().required('Required'),
  secondaryDelivery: AdditionalAddressSchema,
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

    case SHIPMENT_OPTIONS.NTS:
      switch (userRole) {
        case roleTypes.TOO: {
          return {
            schema: ntsShipmentTOOSchema,
            showPickupFields: true,
            showDeliveryFields: false,
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
            showPickupFields: false,
            showDeliveryFields: true,
          };
        }

        default: {
          throw new Error('unrecognized user role type');
        }
      }

    default:
      throw new Error('unrecognized move type');
  }
}

export default getShipmentOptions;

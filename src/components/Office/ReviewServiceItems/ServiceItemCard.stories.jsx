import React from 'react';

import testParams from '../ServiceItemCalculations/serviceItemTestParams';

import ServiceItemCard from './ServiceItemCard';

import { SHIPMENT_OPTIONS, PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';
import { serviceItemCodes } from 'content/serviceItems';

export default {
  title: 'Office Components/ReviewServiceItems/ServiceItemCards',
  component: ServiceItemCard,
  argTypes: {
    patchPaymentServiceItem: {
      action: 'update status',
    },
  },
};

export const Basic = (args) => (
  <ServiceItemCard
    mtoServiceItemName={serviceItemCodes.CS}
    amount={999.99}
    patchPaymentServiceItem={args.patchPaymentServiceItem}
  />
);

export const HHG = (args) => (
  <ServiceItemCard
    mtoShipmentType={SHIPMENT_OPTIONS.HHG}
    mtoShipmentDepartureDate="04 May 2021"
    mtoShipmentPickupAddress="Fairfield, CA 94535"
    mtoShipmentDestinationAddress="Beverly Hills, CA 90210"
    mtoServiceItemCode="FSC"
    mtoServiceItemName={serviceItemCodes.FSC}
    status={PAYMENT_SERVICE_ITEM_STATUS.REQUESTED}
    paymentServiceItemParams={testParams.FuelSurchage}
    amount={999.99}
    patchPaymentServiceItem={args.patchPaymentServiceItem}
  />
);

export const NTS = (args) => (
  <ServiceItemCard
    mtoShipmentType={SHIPMENT_OPTIONS.NTS}
    mtoShipmentDepartureDate="04 May 2021"
    mtoShipmentPickupAddress="Fairfield, CA 94535"
    mtoShipmentDestinationAddress="Beverly Hills, CA 90210"
    mtoServiceItemCode="FSC"
    mtoServiceItemName={serviceItemCodes.FSC}
    status={PAYMENT_SERVICE_ITEM_STATUS.REQUESTED}
    paymentServiceItemParams={testParams.FuelSurchage}
    amount={999.99}
    patchPaymentServiceItem={args.patchPaymentServiceItem}
  />
);

export const NTSR = (args) => (
  <ServiceItemCard
    mtoShipmentType={SHIPMENT_OPTIONS.NTSR}
    mtoShipmentDepartureDate="04 May 2021"
    mtoShipmentPickupAddress="Fairfield, CA 94535"
    mtoShipmentDestinationAddress="Beverly Hills, CA 90210"
    mtoServiceItemCode="FSC"
    mtoServiceItemName={serviceItemCodes.FSC}
    status={PAYMENT_SERVICE_ITEM_STATUS.REQUESTED}
    paymentServiceItemParams={testParams.FuelSurchage}
    amount={999.99}
    patchPaymentServiceItem={args.patchPaymentServiceItem}
  />
);

export const HHGLonghaulDomestic = (args) => (
  <ServiceItemCard
    mtoShipmentType={SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC}
    mtoShipmentDepartureDate="04 May 2021"
    mtoShipmentPickupAddress="Fairfield, CA 94535"
    mtoShipmentDestinationAddress="Beverly Hills, CA 90210"
    mtoServiceItemCode="FSC"
    mtoServiceItemName={serviceItemCodes.FSC}
    status={PAYMENT_SERVICE_ITEM_STATUS.REQUESTED}
    paymentServiceItemParams={testParams.FuelSurchage}
    amount={999.99}
    patchPaymentServiceItem={args.patchPaymentServiceItem}
  />
);

export const HHGShorthaulDomestic = (args) => (
  <ServiceItemCard
    mtoShipmentType={SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC}
    mtoShipmentDepartureDate="04 May 2021"
    mtoShipmentPickupAddress="Fairfield, CA 94535"
    mtoShipmentDestinationAddress="Beverly Hills, CA 90210"
    mtoServiceItemCode="FSC"
    mtoServiceItemName={serviceItemCodes.FSC}
    status={PAYMENT_SERVICE_ITEM_STATUS.REQUESTED}
    paymentServiceItemParams={testParams.FuelSurchage}
    amount={999.99}
    patchPaymentServiceItem={args.patchPaymentServiceItem}
  />
);

export const NeedsReviewRequestCalculations = (args) => (
  <ServiceItemCard
    mtoShipmentType={SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC}
    mtoShipmentDepartureDate="04 May 2021"
    mtoShipmentPickupAddress="Fairfield, CA 94535"
    mtoShipmentDestinationAddress="Beverly Hills, CA 90210"
    mtoServiceItemCode="FSC"
    mtoServiceItemName={serviceItemCodes.FSC}
    status={PAYMENT_SERVICE_ITEM_STATUS.REQUESTED}
    paymentServiceItemParams={testParams.FuelSurchage}
    amount={999.99}
    patchPaymentServiceItem={args.patchPaymentServiceItem}
  />
);

export const AcceptedRequestComplete = () => (
  <ServiceItemCard
    mtoShipmentType={SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC}
    mtoShipmentDepartureDate="04 May 2021"
    mtoShipmentPickupAddress="Fairfield, CA 94535"
    mtoShipmentDestinationAddress="Beverly Hills, CA 90210"
    mtoServiceItemCode="FSC"
    mtoServiceItemName={serviceItemCodes.FSC}
    status={PAYMENT_SERVICE_ITEM_STATUS.APPROVED}
    paymentServiceItemParams={testParams.FuelSurchage}
    amount={999.99}
    requestComplete
  />
);

export const RejectedRequestComplete = () => (
  <ServiceItemCard
    mtoShipmentType={SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC}
    mtoShipmentDepartureDate="04 May 2021"
    mtoShipmentPickupAddress="Fairfield, CA 94535"
    mtoShipmentDestinationAddress="Beverly Hills, CA 90210"
    mtoServiceItemCode="FSC"
    mtoServiceItemName={serviceItemCodes.FSC}
    status={PAYMENT_SERVICE_ITEM_STATUS.DENIED}
    paymentServiceItemParams={testParams.FuelSurchage}
    rejectionReason="Services were provided by the government"
    amount={999.99}
    requestComplete
  />
);

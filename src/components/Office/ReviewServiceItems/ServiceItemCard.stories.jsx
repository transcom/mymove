import React from 'react';
import { expect } from '@storybook/jest';
import { within, userEvent } from '@storybook/testing-library';

import testParams from '../ServiceItemCalculations/serviceItemTestParams';

import ServiceItemCard from './ServiceItemCard';

import { SHIPMENT_OPTIONS, PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';
import { serviceItemCodes } from 'content/serviceItems';
import { shipmentModificationTypes } from 'constants/shipments';
import { SERVICE_ITEM_CODES } from 'constants/serviceItems';

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

export const EmptyRejectionReasonError = (args) => (
  <ServiceItemCard
    mtoServiceItemName={serviceItemCodes.CS}
    amount={999.99}
    patchPaymentServiceItem={args.patchPaymentServiceItem}
    status={PAYMENT_SERVICE_ITEM_STATUS.DENIED}
  />
);

EmptyRejectionReasonError.play = async ({ canvasElement }) => {
  const canvas = within(canvasElement);

  await expect(canvas.getByRole('textbox', { name: 'Reason for rejection' })).toBeInTheDocument();
  await expect(canvas.getByText('Reject')).toBeInTheDocument();

  // type, then clear, then blur
  await userEvent.type(canvas.getByRole('textbox', { name: 'Reason for rejection' }), 'a{backspace}');
  await userEvent.click(canvas.getByText('Reject'), undefined, { skipPointerEventsCheck: true });
};

export const HHG = (args) => (
  <ServiceItemCard
    mtoShipmentType={SHIPMENT_OPTIONS.HHG}
    mtoShipmentDepartureDate="2020-03-16"
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
    mtoShipmentDepartureDate="2020-03-16"
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
    mtoShipmentDepartureDate="2020-03-16"
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
    mtoShipmentType={SHIPMENT_OPTIONS.HHG}
    mtoShipmentDepartureDate="2020-03-16"
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
    mtoShipmentType={SHIPMENT_OPTIONS.HHG}
    mtoShipmentDepartureDate="2020-03-16"
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

export const HHGCanceled = (args) => (
  <ServiceItemCard
    mtoShipmentType={SHIPMENT_OPTIONS.HHG}
    mtoShipmentDepartureDate="04 May 2021"
    mtoShipmentPickupAddress="Fairfield, CA 94535"
    mtoShipmentDestinationAddress="Beverly Hills, CA 90210"
    mtoServiceItemCode="FSC"
    mtoServiceItemName={serviceItemCodes.FSC}
    mtoShipmentModificationType={shipmentModificationTypes.CANCELED}
    status={PAYMENT_SERVICE_ITEM_STATUS.REQUESTED}
    paymentServiceItemParams={testParams.FuelSurchage}
    amount={999.99}
    patchPaymentServiceItem={args.patchPaymentServiceItem}
  />
);

export const HHGDiverted = (args) => (
  <ServiceItemCard
    mtoShipmentType={SHIPMENT_OPTIONS.HHG}
    mtoShipmentDepartureDate="04 May 2021"
    mtoShipmentPickupAddress="Fairfield, CA 94535"
    mtoShipmentDestinationAddress="Beverly Hills, CA 90210"
    mtoShipmentModificationType={shipmentModificationTypes.DIVERSION}
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
    mtoShipmentType={SHIPMENT_OPTIONS.HHG}
    mtoShipmentDepartureDate="2020-03-16"
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
    mtoShipmentType={SHIPMENT_OPTIONS.HHG}
    mtoShipmentDepartureDate="2020-03-16"
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
    mtoShipmentType={SHIPMENT_OPTIONS.HHG}
    mtoShipmentDepartureDate="2020-03-16"
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

export const DaysInSITAllowance = () => (
  <ServiceItemCard
    mtoShipmentType={SHIPMENT_OPTIONS.HHG}
    mtoShipmentDepartureDate="2021-05-08"
    mtoShipmentPickupAddress="Fairfield, CA 94535"
    mtoShipmentDestinationAddress="Beverly Hills, CA 90210"
    mtoServiceItemCode={SERVICE_ITEM_CODES.DOASIT}
    mtoServiceItemName={serviceItemCodes.DOASIT}
    paymentServiceItemParams={testParams.DomesticOriginAdditionalSIT} // DaysInSIT would be 60
    amount={999.99}
    shipmentSITBalance={{
      previouslyBilledDays: 30,
      previouslyBilledEndDate: '2021-06-08',
      pendingSITDaysInvoiced: 60,
      pendingBilledEndDate: '2021-08-08',
      totalSITDaysAuthorized: 120,
      totalSITDaysRemaining: 30,
      totalSITEndDate: '2021-09-08',
    }}
  />
);

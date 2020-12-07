import React from 'react';

import ServiceItemCard from './ServiceItemCard';

import { SHIPMENT_OPTIONS, PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';

export default {
  title: 'TOO/TIO Components/ReviewServiceItems/ServiceItemCards',
  component: ServiceItemCard,
};

export const Basic = () => <ServiceItemCard serviceItemName="Counseling services" amount={999.99} />;

export const HHG = () => (
  <ServiceItemCard shipmentType={SHIPMENT_OPTIONS.HHG} serviceItemName="Counseling services" amount={999.99} />
);

export const NTS = () => (
  <ServiceItemCard shipmentType={SHIPMENT_OPTIONS.NTS} serviceItemName="Counseling services" amount={999.99} />
);

export const HHGLonghaulDomestic = () => (
  <ServiceItemCard
    shipmentType={SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC}
    serviceItemName="Counseling services"
    amount={999.99}
  />
);

export const HHGShorthaulDomestic = () => (
  <ServiceItemCard
    shipmentType={SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC}
    serviceItemName="Counseling services"
    amount={999.99}
  />
);

export const AcceptedRequestComplete = () => (
  <ServiceItemCard
    shipmentType={SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC}
    serviceItemName="Counseling services"
    status={PAYMENT_SERVICE_ITEM_STATUS.APPROVED}
    amount={999.99}
    requestComplete
  />
);

export const RejectedRequestComplete = () => (
  <ServiceItemCard
    shipmentType={SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC}
    serviceItemName="Counseling services"
    status={PAYMENT_SERVICE_ITEM_STATUS.DENIED}
    rejectionReason="Services were provided by the government"
    amount={999.99}
    requestComplete
  />
);

import React from 'react';

import ServiceItemCard from './ServiceItemCard';

import { SHIPMENT_OPTIONS, PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';

export default {
  title: 'Office Components/ReviewServiceItems/ServiceItemCards',
  component: ServiceItemCard,
};

export const Basic = () => <ServiceItemCard mtoServiceItemName="Counseling services" amount={999.99} />;

export const HHG = () => (
  <ServiceItemCard mtoShipmentType={SHIPMENT_OPTIONS.HHG} mtoServiceItemName="Counseling services" amount={999.99} />
);

export const NTS = () => (
  <ServiceItemCard mtoShipmentType={SHIPMENT_OPTIONS.NTS} mtoServiceItemName="Counseling services" amount={999.99} />
);

export const HHGLonghaulDomestic = () => (
  <ServiceItemCard
    mtoShipmentType={SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC}
    mtoServiceItemName="Counseling services"
    amount={999.99}
  />
);

export const HHGShorthaulDomestic = () => (
  <ServiceItemCard
    mtoShipmentType={SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC}
    mtoServiceItemName="Counseling services"
    amount={999.99}
  />
);

export const AcceptedRequestComplete = () => (
  <ServiceItemCard
    mtoShipmentType={SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC}
    mtoServiceItemName="Counseling services"
    status={PAYMENT_SERVICE_ITEM_STATUS.APPROVED}
    amount={999.99}
    requestComplete
  />
);

export const RejectedRequestComplete = () => (
  <ServiceItemCard
    mtoShipmentType={SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC}
    mtoServiceItemName="Counseling services"
    status={PAYMENT_SERVICE_ITEM_STATUS.DENIED}
    rejectionReason="Services were provided by the government"
    amount={999.99}
    requestComplete
  />
);

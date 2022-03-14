import React from 'react';

import RequestedShipments from './RequestedShipments';
import {
  shipments,
  ntsExternalVendorShipments,
  ordersInfo,
  allowancesInfo,
  customerInfo,
  agents,
  serviceItems,
  moveTaskOrders,
} from './RequestedShipmentsTestData';

export default {
  title: 'Office Components/RequestedShipments',
};

export const withOneShipment = () => (
  <div className="officeApp">
    <RequestedShipments
      mtoShipments={[shipments[0]]}
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoAgents={agents}
      shipmentsStatus="SUBMITTED"
      moveTaskOrder={moveTaskOrders[0]}
      moveCode="TE5TC0DE"
    />
  </div>
);
export const withOneExternalVendorShipment = () => (
  <div className="officeApp">
    <RequestedShipments
      mtoShipments={[ntsExternalVendorShipments[0]]}
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoAgents={agents}
      shipmentsStatus="SUBMITTED"
      moveTaskOrder={moveTaskOrders[0]}
      moveCode="TE5TC0DE"
    />
  </div>
);

export const withCompletedServicesCounseling = () => (
  <div className="officeApp">
    <RequestedShipments
      mtoShipments={[shipments[0]]}
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoAgents={agents}
      shipmentsStatus="SUBMITTED"
      moveTaskOrder={moveTaskOrders[1]}
      moveCode="TE5TC0DE"
    />
  </div>
);

export const withMultipleShipments = () => (
  <div className="officeApp">
    <RequestedShipments
      mtoShipments={shipments}
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoAgents={agents}
      shipmentsStatus="SUBMITTED"
      moveTaskOrder={moveTaskOrders[0]}
      moveCode="TE5TC0DE"
    />
  </div>
);

export const withOneApprovedShipment = () => (
  <div className="officeApp">
    <RequestedShipments
      mtoShipments={[shipments[0]]}
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoAgents={agents}
      shipmentsStatus="APPROVED"
      mtoServiceItems={serviceItems}
      moveTaskOrder={moveTaskOrders[0]}
      moveCode="TE5TC0DE"
    />
  </div>
);

export const withMultipleApprovedShipments = () => (
  <div className="officeApp">
    <RequestedShipments
      mtoShipments={shipments}
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      mtoAgents={agents}
      shipmentsStatus="APPROVED"
      mtoServiceItems={serviceItems}
      moveTaskOrder={moveTaskOrders[0]}
      moveCode="TE5TC0DE"
    />
  </div>
);

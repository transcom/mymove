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

import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

export default {
  title: 'Office Components/RequestedShipments',
  component: RequestedShipments,
  decorators: [
    (Story, context) => {
      // Dont wrap with permissions for the read only tests
      if (context.name.includes('Read Only')) {
        return <Story />;
      }

      // By default, show component with permissions
      return (
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <Story />
        </MockProviders>
      );
    },
  ],
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

export const withOneShipmentReadOnly = () => (
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
export const withOneExternalVendorShipmentReadOnly = () => (
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

export const withCompletedServicesCounselingReadOnly = () => (
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

export const withMultipleShipmentsReadOnly = () => (
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

export const withOneApprovedShipmentReadOnly = () => (
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

export const withMultipleApprovedShipmentsReadOnly = () => (
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

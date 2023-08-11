import React from 'react';

import { shipments, ordersInfo, customerInfo, serviceItemsMSandCS } from './RequestedShipmentsTestData';
import ApprovedRequestedShipments from './ApprovedRequestedShipments';

import { MockProviders, MockRouterProvider } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

export default {
  title: 'Office Components/ApprovedRequestedShipments',
  component: ApprovedRequestedShipments,
  decorators: [
    (Story, context) => {
      // Don't wrap with permissions for the read only tests
      if (context.name.includes('Read Only')) {
        return (
          <MockRouterProvider>
            <Story />
          </MockRouterProvider>
        );
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

export const WithOneApprovedShipment = () => (
  <div className="officeApp">
    <ApprovedRequestedShipments
      mtoShipments={[shipments[0]]}
      ordersInfo={ordersInfo}
      customerInfo={customerInfo}
      mtoServiceItems={serviceItemsMSandCS}
      moveCode="TE5TC0DE"
    />
  </div>
);

export const WithMultipleApprovedShipments = () => (
  <div className="officeApp">
    <ApprovedRequestedShipments
      mtoShipments={shipments}
      ordersInfo={ordersInfo}
      customerInfo={customerInfo}
      mtoServiceItems={serviceItemsMSandCS}
      moveCode="TE5TC0DE"
    />
  </div>
);

export const WithOneApprovedShipmentReadOnly = () => (
  <div className="officeApp">
    <ApprovedRequestedShipments
      mtoShipments={[shipments[0]]}
      ordersInfo={ordersInfo}
      customerInfo={customerInfo}
      mtoServiceItems={serviceItemsMSandCS}
      moveCode="TE5TC0DE"
    />
  </div>
);

export const WithMultipleApprovedShipmentsReadOnly = () => (
  <div className="officeApp">
    <ApprovedRequestedShipments
      mtoShipments={shipments}
      ordersInfo={ordersInfo}
      customerInfo={customerInfo}
      mtoServiceItems={serviceItemsMSandCS}
      moveCode="TE5TC0DE"
    />
  </div>
);

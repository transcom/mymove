import React from 'react';

import { shipments, ordersInfo, customerInfo, serviceItemsMSandCS } from './RequestedShipmentsTestData';
import ApprovedRequestedShipments from './ApprovedRequestedShipments';

import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

export default {
  title: 'Office Components/ApprovedRequestedShipments',
  component: ApprovedRequestedShipments,
  decorators: [
    (Story, context) => {
      // Don't wrap with permissions for the read only tests
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

export const withOneApprovedShipment = () => (
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

export const withMultipleApprovedShipments = () => (
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

export const withOneApprovedShipmentReadOnly = () => (
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

export const withMultipleApprovedShipmentsReadOnly = () => (
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

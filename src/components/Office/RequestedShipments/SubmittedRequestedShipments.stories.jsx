import React from 'react';
import { Provider } from 'react-redux';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import {
  shipments,
  ntsExternalVendorShipments,
  ordersInfo,
  allowancesInfo,
  customerInfo,
  moveTaskOrders,
  zeroIncentivePPM,
} from './RequestedShipmentsTestData';
import SubmittedRequestedShipments from './SubmittedRequestedShipments';

import { MockProviders, MockRouterProvider } from 'testUtils';
import { permissionTypes } from 'constants/permissions';
import { store } from 'shared/store';

const queryClient = new QueryClient();

export default {
  title: 'Office Components/SubmittedRequestedShipments',
  component: SubmittedRequestedShipments,
  decorators: [
    (Story, context) => {
      // Don't wrap with permissions for the read only tests
      if (context.name.includes('Read Only')) {
        return (
          <QueryClientProvider client={queryClient}>
            <Provider store={store}>
              <MockRouterProvider>
                <Story />
              </MockRouterProvider>
            </Provider>
          </QueryClientProvider>
        );
      }

      // By default, show component with permissions
      return (
        <QueryClientProvider client={queryClient}>
          <Provider store={store}>
            <MockProviders permissions={[permissionTypes.updateShipment]}>
              <Story />
            </MockProviders>
          </Provider>
        </QueryClientProvider>
      );
    },
  ],
};

export const withOneShipment = () => (
  <div className="officeApp">
    <SubmittedRequestedShipments
      allowancesInfo={allowancesInfo}
      moveCode="TE5TC0DE"
      mtoShipments={[shipments[0]]}
      ordersInfo={ordersInfo}
      customerInfo={customerInfo}
      moveTaskOrder={moveTaskOrders[0]}
    />
  </div>
);

export const withOneExternalVendorShipment = () => (
  <div className="officeApp">
    <SubmittedRequestedShipments
      allowancesInfo={allowancesInfo}
      moveCode="TE5TC0DE"
      mtoShipments={[ntsExternalVendorShipments[0]]}
      ordersInfo={ordersInfo}
      moveTaskOrder={moveTaskOrders[0]}
      customerInfo={customerInfo}
    />
  </div>
);

export const withZeroIncentivePPM = () => (
  <div className="officeApp">
    <SubmittedRequestedShipments
      allowancesInfo={allowancesInfo}
      moveCode="TE5TC0DE"
      mtoShipments={[zeroIncentivePPM[0]]}
      ordersInfo={ordersInfo}
      moveTaskOrder={moveTaskOrders[0]}
      customerInfo={customerInfo}
    />
  </div>
);

export const withCompletedServicesCounseling = () => (
  <div className="officeApp">
    <SubmittedRequestedShipments
      mtoShipments={[shipments[0]]}
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      moveTaskOrder={moveTaskOrders[1]}
      moveCode="TE5TC0DE"
    />
  </div>
);

export const withMultipleShipments = () => (
  <div className="officeApp">
    <SubmittedRequestedShipments
      mtoShipments={shipments}
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      moveTaskOrder={moveTaskOrders[0]}
      moveCode="TE5TC0DE"
    />
  </div>
);

export const withOneShipmentReadOnly = () => (
  <div className="officeApp">
    <SubmittedRequestedShipments
      mtoShipments={[shipments[0]]}
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      moveTaskOrder={moveTaskOrders[0]}
      moveCode="TE5TC0DE"
    />
  </div>
);
export const withOneExternalVendorShipmentReadOnly = () => (
  <div className="officeApp">
    <SubmittedRequestedShipments
      mtoShipments={[ntsExternalVendorShipments[0]]}
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      moveTaskOrder={moveTaskOrders[0]}
      moveCode="TE5TC0DE"
    />
  </div>
);

export const withCompletedServicesCounselingReadOnly = () => (
  <div className="officeApp">
    <SubmittedRequestedShipments
      mtoShipments={[shipments[0]]}
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      moveTaskOrder={moveTaskOrders[1]}
      moveCode="TE5TC0DE"
    />
  </div>
);

export const withMultipleShipmentsReadOnly = () => (
  <div className="officeApp">
    <SubmittedRequestedShipments
      mtoShipments={shipments}
      ordersInfo={ordersInfo}
      allowancesInfo={allowancesInfo}
      customerInfo={customerInfo}
      moveTaskOrder={moveTaskOrders[0]}
      moveCode="TE5TC0DE"
    />
  </div>
);

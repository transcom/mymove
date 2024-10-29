/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { Provider } from 'react-redux';

import MtoShipmentForm from './MtoShipmentForm';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { store } from 'shared/store';
import { MockRouterProvider } from 'testUtils';
import { ORDERS_TYPE } from 'constants/orders';

const defaultProps = {
  pageList: ['page1', 'anotherPage/:foo/:bar'],
  pageKey: 'page1',
  showLoggedInUser: () => {},
  newDutyLocationAddress: {
    city: 'Fort Benning',
    state: 'GA',
    postalCode: '31905',
  },
  currentResidence: {
    city: 'Fort Benning',
    state: 'GA',
    postalCode: '31905',
    streetAddress1: '123 Main',
    streetAddress2: '',
  },
  useCurrentResidence: false,
  mtoShipment: {
    destinationAddress: undefined,
  },
  orders: {
    orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
    authorizedWeight: 5000,
    entitlement: {
      proGear: 1000,
      proGearSpouse: 100,
    },
  },
  isCreatePage: true,
};

const mockMtoShipment = {
  id: 'mock id',
  moveTaskOrderId: 'mock move id',
  customerRemarks: 'mock remarks',
  requestedPickupDate: '2020-03-01',
  requestedDeliveryDate: '2020-03-30',
  agents: [
    {
      firstName: 'mock receiving',
      lastName: 'agent',
      telephone: '2225551234',
      email: 'mock.delivery.agent@example.com',
      agentType: 'RECEIVING_AGENT',
    },
    {
      firstName: 'Mock Releasing',
      lastName: 'Agent Jr, PhD, MD, DDS',
      telephone: '3335551234',
      email: 'mock.pickup.agent@example.com',
      agentType: 'RELEASING_AGENT',
    },
  ],
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
};

export default {
  title: 'Customer Components / Forms / MtoShipmentForm',
};

function renderStory(props) {
  return (
    <Provider store={store}>
      <MockRouterProvider>
        <MtoShipmentForm {...defaultProps} {...props} />
      </MockRouterProvider>
    </Provider>
  );
}

// create shipment stories (form should not prefill customer data)
export const HHGShipment = () => renderStory({ shipmentType: SHIPMENT_OPTIONS.HHG });
export const HHGShipmentRetiree = () =>
  renderStory({ shipmentType: SHIPMENT_OPTIONS.HHG, orders: { orders_type: 'RETIREMENT', authorizedWeight: 5000 } });
export const NTSReleaseShipment = () => renderStory({ shipmentType: SHIPMENT_OPTIONS.NTSR });
export const NTSShipment = () => renderStory({ shipmentType: SHIPMENT_OPTIONS.NTS });
export const UBShipment = () => renderStory({ shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE });

// edit shipment stories (form should prefill)
export const EditHHGShipment = () =>
  renderStory({
    shipmentType: SHIPMENT_OPTIONS.HHG,
    isCreatePage: false,
    mtoShipment: mockMtoShipment,
  });
export const EditNTSReleaseShipment = () =>
  renderStory({
    shipmentType: SHIPMENT_OPTIONS.NTSR,
    isCreatePage: false,
    mtoShipment: mockMtoShipment,
  });
export const EditNTSShipment = () =>
  renderStory({
    shipmentType: SHIPMENT_OPTIONS.NTS,
    isCreatePage: false,
    mtoShipment: mockMtoShipment,
  });
export const EditUBShipment = () =>
  renderStory({
    shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
    isCreatePage: false,
    mtoShipment: mockMtoShipment,
  });

export const EditShipmentAsSeparatee = () =>
  renderStory({
    shipmentType: SHIPMENT_OPTIONS.HHG,
    isCreatePage: false,
    mtoShipment: mockMtoShipment,
    orders: { orders_type: 'SEPARATION', authorizedWeight: 5000 },
  });

export const EditHHGShipmentWithSecondaryAddresses = () => {
  const extendedShipment = {
    ...mockMtoShipment,
    secondaryPickupAddress: {
      streetAddress1: '142 E Barrel Hoop Circle',
      streetAddress2: '#4A',
      city: 'Corpus Christi',
      state: 'TX',
      postalCode: '78412',
    },
    secondaryDeliveryAddress: {
      streetAddress1: '3373 NW Martin Luther King Jr Blvd',
      streetAddress2: '',
      city: mockMtoShipment.destinationAddress.city,
      state: mockMtoShipment.destinationAddress.state,
      postalCode: mockMtoShipment.destinationAddress.postalCode,
    },
  };

  return renderStory({
    shipmentType: SHIPMENT_OPTIONS.HHG,
    isCreatePage: false,
    mtoShipment: extendedShipment,
  });
};

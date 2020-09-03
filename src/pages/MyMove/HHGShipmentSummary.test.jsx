import React from 'react';
import { mount } from 'enzyme';

import HHGShipmentSummary from './HHGShipmentSummary';

import { MockProviders } from 'testUtils';

describe('Review -> Hhg Shipment Summary', () => {
  const mtoShipment = {
    agents: [
      {
        agentType: 'RELEASING_AGENT',
        createdAt: '0001-01-01T00:00:00.000Z',
        email: 'ra@example.com',
        firstName: 'Winnie',
        id: '00000000-0000-0000-0000-000000000000',
        lastName: 'The Pooh',
        mtoShipmentID: '00000000-0000-0000-0000-000000000000',
        phone: '415-444-4444',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
      {
        agentType: 'RECEIVING_AGENT',
        createdAt: '0001-01-01T00:00:00.000Z',
        email: 'bla@example.com',
        firstName: 'Piglet',
        id: '00000000-0000-0000-0000-000000000000',
        lastName: 'Pigleton',
        mtoShipmentID: '00000000-0000-0000-0000-000000000000',
        phone: '415-555-5555',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
    ],
    createdAt: '2020-07-29T00:17:53.236Z',
    customerRemarks: 'please be careful with the move!',
    destinationAddress: {
      city: 'San Francisco',
      id: '0fda108d-6c6c-44c8-b5ae-485b779f7539',
      postal_code: '94611',
      state: 'CA',
      street_address_1: '666 no',
    },
    id: '3dc3c94f-8264-4dd6-85e0-9a0ec1af3433',
    moveTaskOrderID: 'b21536b7-22a3-43c1-a4a7-3a8c392c1ad5',
    pickupAddress: {
      city: 'San Francisco',
      id: '3ea70395-e15d-485b-8b5a-51549069b9f0',
      postal_code: '94611',
      state: 'CA',
      street_address_1: '666 no',
    },
    requestedDeliveryDate: '2020-07-31',
    requestedPickupDate: '2020-07-30',
    shipmentType: 'HHG',
    updatedAt: '2020-07-29T00:17:53.236Z',
  };
  const newDutyStationPostalCode = '94703';
  it('Should display shipment details', () => {
    const wrapper = mount(
      <MockProviders initialEntries={['/']}>
        <HHGShipmentSummary
          mtoShipment={mtoShipment}
          newDutyStationPostalCode={newDutyStationPostalCode}
          movePath="123"
        />
      </MockProviders>,
    );
    // Address
    expect(wrapper.find(HHGShipmentSummary).html()).toContain('666 no');
    // Agent name
    expect(wrapper.find(HHGShipmentSummary).html()).toContain('Winnie The Pooh');
    // Customer Remarks
    expect(wrapper.find(HHGShipmentSummary).html()).toContain('please be careful with the move!');
  });
});

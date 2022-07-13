import React from 'react';
import { render, screen } from '@testing-library/react';

import EvaluationReportShipmentInfo from './EvaluationReportShipmentInfo';

const hhgShipment = {
  actualPickupDate: '2020-03-16',
  createdAt: '2022-07-12T19:38:35.886Z',
  customerRemarks: 'Please treat gently',
  destinationAddress: {
    city: 'Fairfield',
    country: 'US',
    eTag: 'MjAyMi0wNy0xMlQxOTozODozNS44ODQwN1o=',
    id: 'd2aeb8a1-2ddf-4cfc-9067-66a1ad8d5115',
    postalCode: '94535',
    state: 'CA',
    streetAddress1: '987 Any Avenue',
    streetAddress2: 'P.O. Box 9876',
    streetAddress3: 'c/o Some Person',
  },
  eTag: 'MjAyMi0wNy0xMlQxOTozODozNS44ODYyNjFa',
  id: 'c3c64a08-778d-4f9f-8b67-b2502e0fb5e9',
  moveTaskOrderID: 'f256a0fe-5001-46d3-b1ab-a877f75599d4',
  pickupAddress: {
    city: 'Beverly Hills',
    country: 'US',
    eTag: 'MjAyMi0wNy0xMlQxOTozODozNS44ODIyMzJa',
    id: '46dd581a-6eca-44c3-b4e4-be886635f8ab',
    postalCode: '90210',
    state: 'CA',
    streetAddress1: '123 Any Street',
    streetAddress2: 'P.O. Box 12345',
    streetAddress3: 'c/o Some Person',
  },
  primeActualWeight: 980,
  requestedDeliveryDate: '2020-03-15',
  requestedPickupDate: '2020-03-15',
  scheduledPickupDate: '2020-03-16',
  shipmentType: 'HHG',
  status: 'SUBMITTED',
  updatedAt: '2022-07-12T19:38:35.886Z',
};
describe('EvaluationReportShipmentInfo', () => {
  it('renders HHG shipment', () => {
    render(<EvaluationReportShipmentInfo shipment={hhgShipment} />);
    expect(screen.getByRole('heading', { level: 4, name: /HHG/ })).toBeInTheDocument();
    expect(screen.getByText(/123 Any Street/)).toBeInTheDocument();
    expect(screen.getByText(/987 Any Avenue/)).toBeInTheDocument();
  });
});

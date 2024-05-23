import React from 'react';
import { render } from '@testing-library/react';

import '@testing-library/jest-dom/extend-expect';
import ServiceItem from './ServiceItem';

const mockServiceItem = {
  reServiceName: 'Test Service',
  status: 'APPROVED',
  id: '123',
  reServiceCode: 'DPK',
  eTag: 'abc123',
};

const mockSITServiceItem = {
  reServiceName: 'Test Service',
  status: 'APPROVED',
  id: '456',
  reServiceCode: 'DOFSIT',
  eTag: 'abc123',
};

const mockMtoShipment = {
  primeActualWeight: 500,
};

describe('ServiceItem Component', () => {
  it('renders all fields except for Shipment Weight field', () => {
    const { getByText, getByRole, queryByText } = render(
      <ServiceItem serviceItem={mockServiceItem} mtoShipment={mockMtoShipment} />,
    );

    expect(getByRole('heading', { name: 'Test Service' })).toBeInTheDocument();
    expect(getByText('Status:')).toBeInTheDocument();
    expect(getByText('ID:')).toBeInTheDocument();
    expect(getByText('Service Code:')).toBeInTheDocument();
    expect(getByText('Service Name:')).toBeInTheDocument();
    expect(getByText('eTag:')).toBeInTheDocument();
    expect(queryByText('Shipment Weight (pounds):')).toBeNull();
  });

  it('renders shipment weight when service code is a SIT item', () => {
    const { getByText, getByRole } = render(
      <ServiceItem serviceItem={mockSITServiceItem} mtoShipment={mockMtoShipment} />,
    );

    expect(getByRole('heading', { name: 'Test Service' })).toBeInTheDocument();
    expect(getByText('Status:')).toBeInTheDocument();
    expect(getByText('ID:')).toBeInTheDocument();
    expect(getByText('Service Code:')).toBeInTheDocument();
    expect(getByText('Service Name:')).toBeInTheDocument();
    expect(getByText('eTag:')).toBeInTheDocument();
    expect(getByText('Shipment Weight (pounds):')).toBeInTheDocument();
  });
});

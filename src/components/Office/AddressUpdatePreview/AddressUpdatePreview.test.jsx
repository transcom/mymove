import React from 'react';
import { render, screen } from '@testing-library/react';

import AddressUpdatePreview from './AddressUpdatePreview';

const mockDeliveryAddressUpdate = {
  contractorRemarks: 'Test Contractor Remark',
  id: 'c49f7921-5a6e-46b4-bb39-022583574453',
  newAddress: {
    city: 'Beverly Hills',
    country: 'US',
    eTag: 'MjAyMy0wNy0xN1QxODowODowNi42NTU5MTVa',
    id: '6b57ce91-cabd-4e3b-9f48-ed4627d4878f',
    postalCode: '90210',
    state: 'CA',
    streetAddress1: '123 Any Street',
  },
  originalAddress: {
    city: 'Fairfield',
    country: 'US',
    id: '92509013-aafc-4892-a476-2e3b97e6933d',
    postalCode: '94535',
    state: 'CA',
    streetAddress1: '987 Any Avenue',
  },
  shipmentID: '5c84bcf3-92f7-448f-b0e1-e5378b6806df',
  status: 'REQUESTED',
};

describe('AddressUpdatePreview', () => {
  it('renders all of the address preview information', () => {
    render(<AddressUpdatePreview deliveryAddressUpdate={mockDeliveryAddressUpdate} />);

    // Heading and alert present
    expect(screen.getByRole('heading', { name: 'Delivery location' })).toBeInTheDocument();
    expect(screen.getByTestId('alert')).toBeInTheDocument();
    expect(screen.getByTestId('alert')).toHaveTextContent(
      'If approved, the requested update to the delivery location will change one or all of the following:' +
        'Service area.' +
        'Mileage bracket for direct delivery.' +
        'ZIP3 resulting in Domestic Shorthaul (DSH) changing to Domestic Linehaul (DLH) or vice versa.' +
        'Approvals will result in updated pricing for this shipment. Customer may be subject to excess costs.',
    );

    // Address change information
    const addressChangePreview = screen.getByTestId('address-change-preview');
    expect(addressChangePreview).toBeInTheDocument();

    const addresses = screen.getAllByTestId('two-line-address');
    expect(addresses).toHaveLength(2);

    // Original Address
    expect(addressChangePreview).toHaveTextContent('Original delivery location');
    expect(addresses[0]).toHaveTextContent('987 Any Avenue');
    expect(addresses[0]).toHaveTextContent('Fairfield, CA 94535');

    // New Address
    expect(addressChangePreview).toHaveTextContent('Requested delivery location');
    expect(addresses[1]).toHaveTextContent('123 Any Street');
    expect(addresses[1]).toHaveTextContent('Beverly Hills, CA 90210');

    // Request details (contractor remarks)
    expect(addressChangePreview).toHaveTextContent('Update request details');
    expect(addressChangePreview).toHaveTextContent('Contractor remarks: Test Contractor Remark');
  });
});

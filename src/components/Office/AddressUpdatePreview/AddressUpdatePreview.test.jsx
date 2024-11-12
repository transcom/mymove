import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';

import AddressUpdatePreview from './AddressUpdatePreview';

const mockDeliveryAddressUpdateWithoutSIT = {
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
const mockDeliveryAddressUpdateWithSIT = {
  contractorRemarks: 'hello',
  id: '5b1e566e-de89-4523-897f-16d7723a7a64',
  newAddress: {
    city: 'Fairfield',
    eTag: 'MjAyNC0wMS0yMlQyMDo1MTo1NS4xNTQzMjJa',
    id: 'ad28a8df-0301-4cac-b88f-75b42fc491a7',
    postalCode: '73064',
    state: 'CA',
    streetAddress1: '987 Any Avenue',
    streetAddress2: 'P.O. Box 9876',
    streetAddress3: 'c/o Some Person',
  },
  newSitDistanceBetween: 55,
  oldSitDistanceBetween: 1,
  originalAddress: {
    city: 'Fairfield',
    country: 'US',
    eTag: 'MjAyNC0wMS0wM1QyMToyODoyOS4zNDUxNzFa',
    id: 'ac8654e1-8a31-45ed-991b-8c28222cf877',
    postalCode: '94535',
    state: 'CA',
    streetAddress1: '987 Any Avenue',
    streetAddress2: 'P.O. Box 9876',
    streetAddress3: 'c/o Some Person',
  },
  shipmentID: 'fde0a71f-b984-4abf-8491-2e51f41ab1b9',
  sitOriginalAddress: {
    city: 'Fairfield',
    country: 'US',
    eTag: 'MjAyNC0wMS0wM1QyMToyODozMi4wMzg5MTda',
    id: 'acea509d-0add-4bde-95d6-c4f10247f9d3',
    postalCode: '94535',
    state: 'CA',
    streetAddress1: '987 Any Avenue',
    streetAddress2: 'P.O. Box 9876',
    streetAddress3: 'c/o Some Person',
  },
  status: 'REQUESTED',
};
describe('AddressUpdatePreview', () => {
  it('renders all of the address preview information', async () => {
    render(<AddressUpdatePreview deliveryAddressUpdate={mockDeliveryAddressUpdateWithoutSIT} />);
    // Heading and alert present
    expect(screen.getByRole('heading', { name: 'Delivery Address' })).toBeInTheDocument();
    expect(screen.getByTestId('alert')).toBeInTheDocument();
    expect(screen.getByTestId('alert')).toHaveTextContent(
      'If approved, the requested update to the delivery address will change one or all of the following:' +
        'Service area.' +
        'Mileage bracket for direct delivery.' +
        'ZIP3 resulting in Domestic Shorthaul (DSH) changing to Domestic Linehaul (DLH) or vice versa.' +
        'Approvals will result in updated pricing for this shipment. Customer may be subject to excess costs.',
    );
    // since there are no destination service items in this shipment, this alert should not show up
    expect(screen.queryByTestId('destSitAlert')).toBeNull();
    // Address change information
    const addressChangePreview = screen.getByTestId('address-change-preview');
    expect(addressChangePreview).toBeInTheDocument();
    const addresses = screen.getAllByTestId('two-line-address');
    expect(addresses).toHaveLength(2);
    // Original Address
    expect(addressChangePreview).toHaveTextContent('Original delivery address');
    expect(addresses[0]).toHaveTextContent('987 Any Avenue');
    expect(addresses[0]).toHaveTextContent('Fairfield, CA 94535');
    // New Address
    expect(addressChangePreview).toHaveTextContent('Requested delivery address');
    expect(addresses[1]).toHaveTextContent('123 Any Street');
    expect(addresses[1]).toHaveTextContent('Beverly Hills, CA 90210');
    // Request details (contractor remarks)
    expect(addressChangePreview).toHaveTextContent('Update request details');
    expect(addressChangePreview).toHaveTextContent('Contractor remarks: Test Contractor Remark');
    // if the delivery address update doesn't have data, then this will be falsy
    await waitFor(() => {
      expect(screen.queryByTestId('destSitAlert')).not.toBeInTheDocument();
    });
  });
  it('renders the destination SIT alert when shipment contains dest SIT service items', () => {
    render(<AddressUpdatePreview deliveryAddressUpdate={mockDeliveryAddressUpdateWithSIT} />);
    // Heading and alert present
    expect(screen.getByRole('heading', { name: 'Delivery Address' })).toBeInTheDocument();
    expect(screen.getByTestId('destSitAlert')).toBeInTheDocument();
    expect(screen.getByTestId('destSitAlert')).toHaveTextContent(
      'Approval of this address change request will result in SIT Delivery > 50 Miles.' +
        'Updated Mileage for SIT: 55 miles',
    );
  });
});

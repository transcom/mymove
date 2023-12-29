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

const destSitServiceItemsNone = [];

const destSitServiceItemsSeveral = [
  {
    approvedAt: '2023-12-29T15:31:57.041Z',
    createdAt: '2023-12-29T15:27:55.909Z',
    deletedAt: '0001-01-01',
    eTag: 'MjAyMy0xMi0yOVQxNTozMTo1Ny4wNTUxMTNa',
    id: '447c4919-3311-4d3d-9067-a5585ba692ad',
    moveTaskOrderID: 'aa8dfe13-266a-4956-ac60-01c2355c06d3',
    mtoShipmentID: 'be3349f4-333d-4633-8708-d9c1147cd407',
    reServiceCode: 'DDASIT',
    reServiceID: 'a0ead168-7469-4cb6-bc5b-2ebef5a38f92',
    reServiceName: "Domestic destination add'l SIT",
    reason: 'LET THE PEOPLE KNOW',
    sitDepartureDate: '2024-01-06T00:00:00.000Z',
    sitEntryDate: '2024-01-05T00:00:00.000Z',
    status: 'APPROVED',
    submittedAt: '0001-01-01',
    updatedAt: '0001-01-01T00:00:00.000Z',
  },
  {
    approvedAt: '2023-12-29T15:31:57.912Z',
    createdAt: '2023-12-29T15:27:55.920Z',
    deletedAt: '0001-01-01',
    eTag: 'MjAyMy0xMi0yOVQxNTozMTo1Ny45MjA3Njla',
    id: '0163ae1a-d6b8-468d-9ec5-49f289796819',
    moveTaskOrderID: 'aa8dfe13-266a-4956-ac60-01c2355c06d3',
    mtoShipmentID: 'be3349f4-333d-4633-8708-d9c1147cd407',
    reServiceCode: 'DDDSIT',
    reServiceID: '5c80f3b5-548e-4077-9b8e-8d0390e73668',
    reServiceName: 'Domestic destination SIT delivery',
    reason: 'LET THE PEOPLE KNOW',
    sitDepartureDate: '2024-01-06T00:00:00.000Z',
    sitEntryDate: '2024-01-05T00:00:00.000Z',
    status: 'APPROVED',
    submittedAt: '0001-01-01',
    updatedAt: '0001-01-01T00:00:00.000Z',
  },
  {
    approvedAt: '2023-12-29T15:31:58.538Z',
    createdAt: '2023-12-29T15:27:55.928Z',
    deletedAt: '0001-01-01',
    eTag: 'MjAyMy0xMi0yOVQxNTozMTo1OC41NDQ0NTJa',
    id: 'b582cc4c-23ae-4529-be20-608b305d9dc7',
    moveTaskOrderID: 'aa8dfe13-266a-4956-ac60-01c2355c06d3',
    mtoShipmentID: 'be3349f4-333d-4633-8708-d9c1147cd407',
    reServiceCode: 'DDSFSC',
    reServiceID: 'b208e0af-3176-4c8a-97ea-bd247c18f43d',
    reServiceName: 'Domestic destination SIT fuel surcharge',
    reason: 'LET THE PEOPLE KNOW',
    sitDepartureDate: '2024-01-06T00:00:00.000Z',
    sitEntryDate: '2024-01-05T00:00:00.000Z',
    status: 'APPROVED',
    submittedAt: '0001-01-01',
    updatedAt: '0001-01-01T00:00:00.000Z',
  },
  {
    approvedAt: '2023-12-29T15:31:59.239Z',
    createdAt: '2023-12-29T15:27:55.837Z',
    deletedAt: '0001-01-01',
    eTag: 'MjAyMy0xMi0yOVQxNTozMTo1OS4yNDU0MDRa',
    id: 'a69e8cb9-5e46-43a5-92e6-27f1f073d92e',
    moveTaskOrderID: 'aa8dfe13-266a-4956-ac60-01c2355c06d3',
    mtoShipmentID: 'be3349f4-333d-4633-8708-d9c1147cd407',
    reServiceCode: 'DDFSIT',
    reServiceID: 'd0561c49-e1a9-40b8-a739-3e639a9d77af',
    reServiceName: 'Domestic destination 1st day SIT',
    reason: 'LET THE PEOPLE KNOW',
    sitDepartureDate: '2024-01-06T00:00:00.000Z',
    sitEntryDate: '2024-01-05T00:00:00.000Z',
    status: 'APPROVED',
    submittedAt: '0001-01-01',
    updatedAt: '0001-01-01T00:00:00.000Z',
  },
];

describe('AddressUpdatePreview', () => {
  it('renders all of the address preview information', () => {
    render(
      <AddressUpdatePreview
        deliveryAddressUpdate={mockDeliveryAddressUpdate}
        destSitServiceItems={destSitServiceItemsNone}
      />,
    );

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

    // since there are no destination service items in this shipment, this alert should not show up
    expect(screen.queryByTestId('destSitAlert')).toBeNull();

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

  it('renders the destination SIT alert when shipment contains dest SIT service items', () => {
    render(
      <AddressUpdatePreview
        deliveryAddressUpdate={mockDeliveryAddressUpdate}
        destSitServiceItems={destSitServiceItemsSeveral}
      />,
    );

    // Heading and alert present
    expect(screen.getByRole('heading', { name: 'Delivery location' })).toBeInTheDocument();
    expect(screen.getByTestId('destSitAlert')).toBeInTheDocument();
    expect(screen.getByTestId('destSitAlert')).toHaveTextContent(
      'This shipment contains 4 destination SIT service items. If approved, this could change the following: ' +
        'SIT Delivery out over 50 miles or under 50 miles' +
        'Approvals will result in updated pricing for the service item. Customer may be subject to excess costs.',
    );
  });
});

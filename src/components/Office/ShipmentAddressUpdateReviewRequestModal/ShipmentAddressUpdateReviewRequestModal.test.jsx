import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';

import { ShipmentAddressUpdateReviewRequestModal } from './ShipmentAddressUpdateReviewRequestModal';

import { shipmentTypes } from 'constants/shipments';

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

const mockOnClose = jest.fn();

afterEach(() => {
  jest.resetAllMocks();
});

describe('ShipmentAddressUpdateReviewRequestModal', () => {
  it('renders the modal', async () => {
    render(
      <ShipmentAddressUpdateReviewRequestModal
        shipmentType={shipmentTypes.HHG}
        deliveryAddressUpdate={mockDeliveryAddressUpdate}
        onClose={mockOnClose}
      />,
    );

    // console.log(screen.debug());
    await waitFor(() => {
      // Shipment type flag
      expect(screen.getByTestId('tag')).toHaveTextContent('HHG');

      // Heading
      expect(screen.getByRole('heading', { level: 2, name: 'Review request' })).toBeInTheDocument();

      // Address update preview component
      expect(screen.getByRole('heading', { level: 3, name: 'Delivery location' })).toBeInTheDocument();

      // Form
      expect(screen.getByRole('heading', { level: 4, name: 'Review Request' })).toBeInTheDocument();

      // Form fields
      expect(screen.getByText('Approve address change?')).toBeInTheDocument();
      expect(screen.getByRole('radio', { name: 'Yes' })).toBeInTheDocument();
      expect(screen.getByRole('radio', { name: 'No' })).toBeInTheDocument();

      expect(screen.getByLabelText('Office remarks')).toBeInTheDocument();
      expect(screen.getByText('Office remarks will be sent to the contractor.')).toBeInTheDocument();
      expect(screen.getByTestId('officeRemarks')).toBeInTheDocument();

      // Form buttons (Save and Cancel)
      expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    });
  });
});

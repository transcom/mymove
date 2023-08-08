import React from 'react';
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { v4 as uuidv4 } from 'uuid';

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

const mockShipment = {
  id: uuidv4(),
  shipmentType: shipmentTypes.HHG,
  deliveryAddressUpdate: mockDeliveryAddressUpdate,
  eTag: 'eTag',
};

afterEach(() => {
  jest.resetAllMocks();
});

describe('ShipmentAddressUpdateReviewRequestModal', () => {
  it('renders the modal', async () => {
    render(
      <ShipmentAddressUpdateReviewRequestModal shipment={mockShipment} onSubmit={jest.fn()} onClose={jest.fn()} />,
    );

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
      const approvalQuestion = screen.getByRole('group', { name: 'Approve address change?' });
      expect(approvalQuestion).toBeInTheDocument();
      const approvalYes = within(approvalQuestion).getByRole('radio', { name: 'Yes' });
      const approvalNo = within(approvalQuestion).getByRole('radio', { name: 'No' });
      expect(approvalYes).toBeInTheDocument();
      expect(approvalNo).toBeInTheDocument();

      expect(screen.getByLabelText('Office remarks')).toBeInTheDocument();
      expect(screen.getByText('Office remarks will be sent to the contractor.')).toBeInTheDocument();
      expect(screen.getByTestId('officeRemarks')).toBeInTheDocument();

      // Form buttons (Save and Cancel)
      expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    });
  });

  it('Runs an onClose callback on close', async () => {
    const user = userEvent.setup();

    const onClose = jest.fn();

    render(<ShipmentAddressUpdateReviewRequestModal shipment={mockShipment} onSubmit={jest.fn()} onClose={onClose} />);

    const cancel = screen.getByRole('button', { name: 'Cancel' });

    expect(cancel).toBeInTheDocument();

    await user.click(cancel);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('Runs an onSubmit callback on save', async () => {
    const user = userEvent.setup();

    const onSubmit = jest.fn();
    const officeRemarksAnswer = 'Here are my remarks from the office';

    render(<ShipmentAddressUpdateReviewRequestModal shipment={mockShipment} onSubmit={onSubmit} onClose={jest.fn()} />);

    const approvalQuestion = screen.getByRole('group', { name: 'Approve address change?' });
    const approvalYes = within(approvalQuestion).getByRole('radio', { name: 'Yes' });
    const officeRemarks = screen.getByLabelText('Office remarks');
    const save = screen.getByRole('button', { name: 'Save' });

    await user.click(approvalYes);
    await user.type(officeRemarks, officeRemarksAnswer);

    expect(approvalYes).toBeChecked();
    expect(officeRemarks).toHaveValue(officeRemarksAnswer);

    await user.click(save);

    expect(onSubmit).toHaveBeenCalledTimes(1);
    expect(onSubmit).toHaveBeenCalledWith(mockShipment.id, mockShipment.eTag, 'APPROVED', officeRemarksAnswer);
  });
});

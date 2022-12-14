import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';

import ReviewProGear from './ReviewProGear';

import ppmDocumentStatus from 'constants/ppms';

beforeEach(() => {
  jest.clearAllMocks();
});

const defaultProps = {
  ppmShipment: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    actualMoveDate: '02-Dec-22',
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
    hasReceivedAdvance: true,
    advanceAmountReceived: 60000,
  },
  tripNumber: 1,
  ppmNumber: 1,
};

const proGearRequiredProps = {
  proGear: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    selfProGear: true,
    proGearWeight: 400,
    description: 'Kia Forte',
    missingWeightTicket: false,
  },
};

const missingWeightTicketProps = {
  proGear: {
    ...proGearRequiredProps.proGear,
    missingWeightTicket: true,
  },
};

const rejectedProps = {
  proGear: {
    ...proGearRequiredProps.proGear,
    status: ppmDocumentStatus.REJECTED,
    reason: 'Rejection reason',
  },
};

describe('ReviewProGear component', () => {
  describe('displays form', () => {
    it('renders blank form on load with defaults', async () => {
      render(<ReviewProGear {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 3, name: 'Pro-gear 1' })).toBeInTheDocument();
      });
      expect(screen.getByText('Belongs to')).toBeInTheDocument();
      expect(screen.getByLabelText('Customer')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Spouse')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Description')).toBeInTheDocument();

      expect(screen.getByText('Pro-gear weight')).toBeInTheDocument();
      expect(screen.getByLabelText('Weight tickets')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Constructed weight')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByLabelText(/Shipment's pro-gear weight/)).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByRole('heading', { level: 3, name: 'Review pro-gear 1' })).toBeInTheDocument();
      expect(screen.getByText('Add a review for this pro-gear')).toBeInTheDocument();
      expect(screen.getByLabelText('Approve')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Reject')).toBeInstanceOf(HTMLInputElement);
    });

    it('populates edit form with existing weight ticket values', async () => {
      render(<ReviewProGear {...defaultProps} {...proGearRequiredProps} />);

      await waitFor(() => {
        expect(screen.getByLabelText('Customer')).toBeChecked();
      });
      expect(screen.getByText('Kia Forte')).toBeInTheDocument();
      expect(screen.getByLabelText(/Shipment's pro-gear weight/)).toHaveDisplayValue('400');
    });

    it('populates edit form when weight ticket is missing', async () => {
      render(<ReviewProGear {...defaultProps} {...missingWeightTicketProps} />);
      await waitFor(() => {
        expect(screen.getByLabelText('Constructed weight')).toBeChecked();
      });
      expect(screen.getByText('Constructed pro-gear weight')).toBeInTheDocument();
    });

    it('displays remaining character count', async () => {
      render(<ReviewProGear {...defaultProps} {...rejectedProps} />);
      await waitFor(() => {
        expect(screen.getByLabelText('Reason')).toHaveDisplayValue('Rejection reason');
      });
      expect(screen.getByText('500 characters')).toBeInTheDocument();
    });
  });
});

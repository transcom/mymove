import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';

import ReviewWeightTicket from './ReviewWeightTicket';

beforeEach(() => {
  jest.clearAllMocks();
});

const defaultProps = {
  ppmShipment: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    actualPickupPostalCode: '90210',
    actualMoveDate: '2022-04-30',
    actualDestinationPostalCode: '94611',
    hasReceivedAdvance: true,
    advanceAmountReceived: 60000,
  },
  tripNumber: 1,
  ppmNumber: 1,
};

const baseWeightTicketProps = {
  id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
  ppmShipmentId: '343bb456-63af-4f76-89bd-7403094a5c4d',
  vehicleDescription: 'Kia Forte',
  emptyWeight: 400,
  fullWeight: 1200,
};

const missingWeightTicketProps = {
  weightTicket: {
    ...baseWeightTicketProps,
    ownsTrailer: false,
    missingEmptyWeightTicket: true,
    missingFullWeightTicket: true,
  },
};

const weightTicketRequiredProps = {
  weightTicket: {
    ...baseWeightTicketProps,
    ownsTrailer: false,
  },
};

const claimableTrailerProps = {
  weightTicket: {
    ...baseWeightTicketProps,
    ownsTrailer: true,
    trailerMeetsCriteria: true,
  },
};

describe('ReviewWeightTicket component', () => {
  describe('displays form', () => {
    it('renders blank form on load with defaults', async () => {
      render(<ReviewWeightTicket {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 3, name: 'Trip 1' })).toBeInTheDocument();
      });

      expect(screen.getByText('Vehicle description')).toBeInTheDocument();
      expect(screen.getByLabelText('Full weight')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByText('Net weight')).toBeInTheDocument();
      expect(screen.getByText('Did they use a trailer they owned?')).toBeInTheDocument();
      expect(screen.getByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('No')).toBeInstanceOf(HTMLInputElement);
      expect(screen.queryByText("Is the trailer's weight claimable?")).not.toBeInTheDocument();

      expect(screen.getByRole('heading', { level: 3, name: 'Review trip 1' })).toBeInTheDocument();
      expect(screen.getByLabelText('Accept')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Reject')).toBeInstanceOf(HTMLInputElement);
    });

    it('populates edit form with existing weight ticket values', async () => {
      render(<ReviewWeightTicket {...defaultProps} {...weightTicketRequiredProps} />);

      await waitFor(() => {
        expect(screen.getByText('Kia Forte')).toBeInTheDocument();
      });
      expect(screen.getByLabelText('Empty weight', { description: 'Weight tickets' })).toHaveDisplayValue('400');
      expect(screen.getByLabelText('Full weight', { description: 'Weight tickets' })).toHaveDisplayValue('1,200');
      expect(screen.getByText('800 lbs')).toBeInTheDocument();
      expect(screen.getByLabelText('No')).toBeChecked();
    });

    it('populates edit form when weight ticket is missing', async () => {
      render(<ReviewWeightTicket {...defaultProps} {...missingWeightTicketProps} />);
      await waitFor(() => {
        expect(screen.getByLabelText('Empty weight', { description: 'Constructed weight' })).toBeInTheDocument();
      });
      expect(screen.getByLabelText('Full weight', { description: 'Constructed weight' })).toBeInTheDocument();
    });

    it('notifies the user when a trailer is claimable, and disables approval', async () => {
      render(<ReviewWeightTicket {...defaultProps} {...claimableTrailerProps} />);
      await waitFor(() => {
        expect(screen.queryByText("Is the trailer's weight claimable?")).toBeInTheDocument();
      });
      expect(screen.queryByText('Proof of ownership is needed to accept this item.')).toBeInTheDocument();
      expect(screen.getByLabelText('Accept')).toHaveAttribute('disabled');
    });
  });
});

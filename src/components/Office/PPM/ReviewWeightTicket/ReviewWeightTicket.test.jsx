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

const missingWeightTicketProps = {
  weightTicket: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    ppmShipmentId: '343bb456-63af-4f76-89bd-7403094a5c4d',
    vehicleDescription: 'Kia Forte',
    emptyWeight: 400,
    fullWeight: 1200,
    ownsTrailer: false,
    missingEmptyWeightTicket: true,
  },
};

const weightTicketRequiredProps = {
  weightTicket: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    ppmShipmentId: '343bb456-63af-4f76-89bd-7403094a5c4d',
    vehicleDescription: 'Kia Forte',
    emptyWeight: 400,
    fullWeight: 1200,
    ownsTrailer: false,
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
      expect(screen.getByText('Weight type')).toBeInTheDocument();
      expect(screen.getByLabelText('Weight tickets')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Constructed weight')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Empty weight')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Full weight')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByText('Net weight')).toBeInTheDocument();
      expect(screen.getByText('Did they use a trailer they owned?')).toBeInTheDocument();
      expect(screen.getByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('No')).toBeInstanceOf(HTMLInputElement);
      expect(screen.queryByText("Is the trailer's weight claimable?")).not.toBeInTheDocument();

      expect(screen.getByRole('heading', { level: 3, name: 'Review trip 1' })).toBeInTheDocument();
      expect(screen.getByLabelText('Approve')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Reject')).toBeInstanceOf(HTMLInputElement);
    });

    it('populates edit form with existing weight ticket values', async () => {
      render(<ReviewWeightTicket {...defaultProps} {...weightTicketRequiredProps} />);

      await waitFor(() => {
        expect(screen.getByText('Kia Forte')).toBeInTheDocument();
      });
      expect(screen.getByLabelText('Empty weight')).toHaveDisplayValue('400');
      expect(screen.getByLabelText('Full weight')).toHaveDisplayValue('1,200');
      expect(screen.getByText('800 lbs')).toBeInTheDocument();
      expect(screen.getByLabelText('No')).toBeChecked();
    });

    it('populates edit form when weight ticket is missing', async () => {
      render(<ReviewWeightTicket {...defaultProps} {...missingWeightTicketProps} />);
      await waitFor(() => {
        expect(screen.getByLabelText('Constructed weight')).toBeChecked();
      });
      expect(screen.getByText('Empty constructed weight')).toBeInTheDocument();
      expect(screen.getByText('Full constructed weight')).toBeInTheDocument();
    });
  });
});

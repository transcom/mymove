import React from 'react';
import { render, waitFor, screen, fireEvent } from '@testing-library/react';
import { act } from 'react-dom/test-utils';

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

const ownsTrailerProps = {
  weightTicket: {
    ...baseWeightTicketProps,
    ownsTrailer: true,
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

    it('toggles the reason field when Reject is selected', async () => {
      render(<ReviewWeightTicket {...defaultProps} {...weightTicketRequiredProps} />);
      await waitFor(() => {
        expect(screen.getByLabelText('Reject')).toBeInstanceOf(HTMLInputElement);
      });
      await act(async () => {
        await fireEvent.click(screen.getByLabelText('Reject'));
      });
      expect(screen.getByLabelText('Reason')).toBeInstanceOf(HTMLTextAreaElement);
      await act(async () => {
        await fireEvent.click(screen.getByLabelText('Accept'));
      });
      expect(screen.queryByLabelText('Reason')).not.toBeInTheDocument();
    });

    it('notifies the user when a trailer is claimable, and disables approval', async () => {
      render(<ReviewWeightTicket {...defaultProps} {...ownsTrailerProps} />);
      await waitFor(() => {
        expect(screen.queryByText("Is the trailer's weight claimable?")).toBeInTheDocument();
      });
      const claimableYesButton = screen.getAllByRole('radio', { name: 'Yes' })[1];
      await act(async () => {
        await fireEvent.click(claimableYesButton);
      });
      expect(screen.queryByText('Proof of ownership is needed to accept this item.')).toBeInTheDocument();
      expect(screen.getByLabelText('Accept')).not.toBeChecked();
      expect(screen.getByLabelText('Accept')).toHaveAttribute('disabled');
    });

    it('notifies the user when a trailer is claimable after toggling ownership', async () => {
      render(<ReviewWeightTicket {...defaultProps} {...claimableTrailerProps} />);
      await waitFor(() => {
        expect(screen.queryByText("Is the trailer's weight claimable?")).toBeInTheDocument();
      });
      const ownedNoButton = screen.getAllByRole('radio', { name: 'No' })[0];
      const ownedYesButton = screen.getAllByRole('radio', { name: 'Yes' })[0];
      await act(async () => {
        await fireEvent.click(ownedNoButton);
        await fireEvent.click(ownedYesButton);
      });
      expect(screen.queryByText('Proof of ownership is needed to accept this item.')).not.toBeInTheDocument();
    });

    it('reenables approval after disabling it and updating weight claimable field', async () => {
      render(<ReviewWeightTicket {...defaultProps} {...claimableTrailerProps} />);
      await waitFor(() => {
        expect(screen.queryByText("Is the trailer's weight claimable?")).toBeInTheDocument();
      });
      const claimableNoButton = screen.getAllByRole('radio', { name: 'No' })[1];
      await act(async () => {
        await fireEvent.click(claimableNoButton);
      });
      expect(screen.queryByText('Proof of ownership is needed to accept this item.')).not.toBeInTheDocument();
      expect(screen.getByLabelText('Accept')).not.toHaveAttribute('disabled');
    });
    describe('shows an error when submitting', () => {
      it('without a status selected', async () => {
        render(<ReviewWeightTicket {...defaultProps} {...claimableTrailerProps} />);
        await waitFor(async () => {
          const form = screen.getByRole('form');
          expect(form).toBeInTheDocument();
          await fireEvent.submit(form);
          expect(screen.getByText('Reviewing this weight ticket is required'));
        });
      });
      it('with Rejected but no reason selected', async () => {
        render(<ReviewWeightTicket {...defaultProps} {...claimableTrailerProps} />);
        await waitFor(async () => {
          const form = screen.getByRole('form');
          expect(form).toBeInTheDocument();
          const rejectionButton = screen.getByTestId('rejectRadio');
          expect(rejectionButton).toBeInTheDocument();
          await fireEvent.click(rejectionButton);
          await fireEvent.submit(form);
          expect(screen.getByText('Add a reason why this weight ticket is rejected'));
        });
      });
    });
  });
});

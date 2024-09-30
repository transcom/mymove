import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import AdvanceForm from 'components/Customer/PPM/Booking/Advance/AdvanceForm';

const defaultProps = {
  onSubmit: jest.fn(),
  onBack: jest.fn(),
  mtoShipment: {
    id: '123',
    ppmShipment: {
      id: '123',
      estimatedIncentive: 1000000,
    },
  },
};

const mtoShipmentProps = {
  onSubmit: jest.fn(),
  onBack: jest.fn(),
  mtoShipment: {
    id: '123',
    ppmShipment: {
      id: '123',
      estimatedIncentive: 1000000,
      hasRequestedAdvance: true,
      advanceAmountRequested: 30000,
    },
  },
};

describe('AdvanceForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load', async () => {
      render(<AdvanceForm {...defaultProps} />);
      expect(
        screen.getByRole('heading', { level: 2, name: 'You can ask for up to $6,000 as an advance' }),
      ).toBeInTheDocument();
      expect(screen.getByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('No')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByText('Would you like to request an advance on your incentive?')).toBeInstanceOf(
        HTMLLegendElement,
      );
    });
  });

  describe('displays conditional inputs', () => {
    it('displays input for amount requested when advance requested is true', async () => {
      render(<AdvanceForm {...defaultProps} />);
      const requestAdvance = screen.getByLabelText('Yes');
      expect(await screen.queryByLabelText(/Amount requested/)).toBeNull();
      expect(
        screen.queryByLabelText(
          "I acknowledge that any advance I'm given will be deducted from my final incentive payment. If my advance ends up being more than my incentive, I will need to repay the difference.",
        ),
      ).toBeNull();
      await userEvent.click(requestAdvance);

      await waitFor(() => {
        expect(screen.getByLabelText(/Amount requested/)).toBeInstanceOf(HTMLInputElement);
        expect(
          screen.getByLabelText(
            "I acknowledge that any advance I'm given will be deducted from my final incentive payment. If my advance ends up being more than my incentive, I will need to repay the difference.",
          ),
        ).toBeInstanceOf(HTMLInputElement);
      });
    });

    it('marks amount requested input as required when conditionally displayed', async () => {
      render(<AdvanceForm {...defaultProps} />);

      const inputHasRequestedAdvance = screen.getByLabelText('Yes');

      await userEvent.click(inputHasRequestedAdvance);

      const advanceAmountRequested = screen.getByLabelText(/Amount requested/);

      await userEvent.click(advanceAmountRequested);
      await userEvent.tab();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();

        const requiredAlerts = screen.getByRole('alert');

        expect(requiredAlerts).toHaveTextContent('Required');
      });
    });

    it('marks amount requested input as min of $1 expected when conditionally displayed', async () => {
      render(<AdvanceForm {...defaultProps} />);

      const inputHasRequestedAdvance = screen.getByLabelText('Yes');

      await userEvent.click(inputHasRequestedAdvance);

      const advanceAmountRequested = screen.getByLabelText(/Amount requested/);

      await userEvent.click(advanceAmountRequested);
      await userEvent.type(advanceAmountRequested, '0');
      await userEvent.tab();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();

        const requiredAlerts = screen.getByRole('alert');

        expect(requiredAlerts).toHaveTextContent(
          "The minimum advance request is $1. If you don't want an advance, select No.",
        );
      });
    });

    it('sets error for requested advance input if it is over allowed amount', async () => {
      render(<AdvanceForm {...defaultProps} />);

      const inputHasRequestedAdvance = screen.getByLabelText('Yes');

      await userEvent.click(inputHasRequestedAdvance);

      const advanceAmountRequested = screen.getByLabelText(/Amount requested/);

      await userEvent.click(advanceAmountRequested);
      await userEvent.type(advanceAmountRequested, '10000');
      await userEvent.tab();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();
      });

      const requiredAlerts = screen.getByRole('alert');

      expect(requiredAlerts).toHaveTextContent('Enter an amount $6,000 or less');
    });
  });

  describe('pull values from the ppm shipment when available', () => {
    it('renders prefilled form on load', async () => {
      render(<AdvanceForm {...mtoShipmentProps} />);
      await waitFor(() => {
        expect(screen.queryByLabelText('Yes').checked).toBe(true);
        expect(screen.getByLabelText(/Amount requested/).value).toBe('300');
      });
    });
  });
});

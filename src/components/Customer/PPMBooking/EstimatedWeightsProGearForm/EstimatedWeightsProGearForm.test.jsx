import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import EstimatedWeightsProGearForm from './EstimatedWeightsProGearForm';

const defaultProps = {
  onSubmit: jest.fn(),
  onBack: jest.fn(),
  orders: {
    has_dependents: true,
  },
  serviceMember: {
    id: '10',
    weight_allotment: {
      total_weight_self: 5000,
      total_weight_self_plus_dependents: 7000,
      pro_gear_weight: 2000,
      pro_gear_weight_spouse: 500,
    },
  },
};

const mtoShipmentProps = {
  ...defaultProps,
  mtoShipment: {
    id: '123',
    ppmShipment: {
      id: '123',
      hasProGear: true,
      proGearWeight: 1000,
      spouseProGearWeight: 100,
      estimatedWeight: 4000,
    },
  },
};

describe('EstimatedWeightsProGearForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load', async () => {
      render(<EstimatedWeightsProGearForm {...defaultProps} />);
      expect(await screen.getByRole('heading', { level: 2, name: 'Full PPM' })).toBeInTheDocument();
      expect(screen.getByRole('heading', { level: 2, name: 'Pro-gear' })).toBeInTheDocument();
      expect(screen.getByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('No')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Estimated weight of this PPM shipment')).toBeInstanceOf(HTMLInputElement);
    });
  });

  describe('displays conditional inputs', () => {
    it('displays secondary pro gear weight inputs when hasProGear is true', async () => {
      render(<EstimatedWeightsProGearForm {...defaultProps} />);
      const hasProGear = await screen.getByLabelText('Yes');
      expect(screen.queryByLabelText('Estimated weight of your pro-gear')).toBeNull();
      expect(screen.queryByLabelText('Estimated weight of your spouse’s pro-gear')).toBeNull();
      userEvent.click(hasProGear);

      await waitFor(() => {
        expect(screen.getByLabelText('Estimated weight of your pro-gear')).toBeInstanceOf(HTMLInputElement);
        expect(screen.getByLabelText('Estimated weight of your spouse’s pro-gear')).toBeInstanceOf(HTMLInputElement);
      });
    });
  });

  describe('pull values from the ppm shipment when available', () => {
    it('renders blank form on load', async () => {
      render(<EstimatedWeightsProGearForm {...mtoShipmentProps} />);

      await waitFor(() => {
        expect(screen.getByLabelText('Estimated weight of this PPM shipment').value).toBe('4,000');
      });
      expect(screen.getByLabelText('Yes').value).toBe('true');
      expect(screen.getByLabelText('Estimated weight of your pro-gear').value).toBe('1,000');
      expect(screen.getByLabelText('Estimated weight of your spouse’s pro-gear').value).toBe('100');
    });
  });

  describe('validates form fields and displays error messages', () => {
    it('marks required inputs when left empty', async () => {
      render(<EstimatedWeightsProGearForm {...defaultProps} />);

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();

        const requiredAlerts = screen.getAllByRole('alert');

        // Estimated PPM Weight
        expect(requiredAlerts[0]).toHaveTextContent('Required');
      });
    });

    it('marks secondary pro gear inputs as required when conditionally displayed', async () => {
      render(<EstimatedWeightsProGearForm {...defaultProps} />);

      const inputHasProGear = screen.getByLabelText('Yes');

      await userEvent.click(inputHasProGear);

      const selfProGear = screen.getByLabelText('Estimated weight of your pro-gear');

      await userEvent.click(selfProGear);
      await userEvent.tab();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();

        const requiredAlerts = screen.getByRole('alert');

        expect(requiredAlerts).toHaveTextContent(
          "Enter a weight into at least one pro-gear field. If you won't have pro-gear, select No above.",
        );
      });
    });
  });
});

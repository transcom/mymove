import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import EstimatedWeightsProGearForm from 'components/Customer/PPM/Booking/EstimatedWeightsProGearForm/EstimatedWeightsProGearForm';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

const defaultProps = {
  onSubmit: jest.fn(),
  onBack: jest.fn(),
  orders: {
    has_dependents: false,
    authorizedWeight: 5000,
    entitlement: {
      proGearWeight: 2000,
      spouseProGearWeight: 500,
    },
  },
  serviceMember: {
    id: '10',
  },
  mtoShipment: {
    id: '123',
    ppmShipment: {
      id: '123',
      sitExpected: false,
      expectedDepartureDate: '2022-12-31',
    },
  },
};

const mtoShipmentProps = {
  ...defaultProps,
  mtoShipment: {
    ...defaultProps.mtoShipment,
    ppmShipment: {
      ...defaultProps.mtoShipment.ppmShipment,
      hasProGear: true,
      proGearWeight: 1000,
      spouseProGearWeight: 100,
      estimatedWeight: 4000,
    },
  },
};

const mtoShipmentPropsWithGunSafe = {
  ...defaultProps,
  mtoShipment: {
    ...defaultProps.mtoShipment,
    ppmShipment: {
      ...defaultProps.mtoShipment.ppmShipment,
      hasProGear: true,
      proGearWeight: 1000,
      spouseProGearWeight: 100,
      estimatedWeight: 4000,
      hasGunSafe: true,
      gunSafeWeight: 500,
    },
  },
};

describe('EstimatedWeightsProGearForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load', async () => {
      render(<EstimatedWeightsProGearForm {...defaultProps} />);
      expect(await screen.getByRole('heading', { level: 2, name: 'PPM' })).toBeInTheDocument();
      expect(screen.getByRole('heading', { level: 2, name: 'Pro-gear' })).toBeInTheDocument();
      expect(screen.getByTestId('hasProGearYes')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTestId('hasProGearNo')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Estimated weight of this PPM shipment/)).toBeInstanceOf(HTMLInputElement);

      expect(screen.queryByText(/Do you have a gun safe that you'll move in this PPM/)).toBeNull();
      expect(screen.queryByTestId('hasGunSafeYes')).not.toBeInTheDocument();
      expect(screen.queryByTestId('hasGunSafeNo')).not.toBeInTheDocument();
    });

    it('renders blank form on load with gun safe (ff on)', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

      await render(<EstimatedWeightsProGearForm {...defaultProps} />);
      expect(await screen.getByRole('heading', { level: 2, name: 'PPM' })).toBeInTheDocument();
      expect(screen.getByRole('heading', { level: 2, name: 'Pro-gear' })).toBeInTheDocument();
      expect(screen.getByTestId('hasProGearYes')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTestId('hasProGearNo')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Estimated weight of this PPM shipment/)).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText(/Do you have a gun safe that you'll move in this PPM/)).toBeInTheDocument();
      expect(screen.getByTestId('hasGunSafeYes')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTestId('hasGunSafeNo')).toBeInstanceOf(HTMLInputElement);
      expect(screen.queryByLabelText(/Estimated weight of your gun safe/)).toBeNull();
    });
  });

  describe('displays conditional inputs', () => {
    it('displays secondary pro gear weight inputs when hasProGear is true', async () => {
      render(<EstimatedWeightsProGearForm {...defaultProps} />);
      const hasProGear = await screen.getByTestId('hasProGearYes');
      expect(screen.queryByLabelText(/Estimated weight of your pro-gear/)).toBeNull();
      expect(screen.queryByLabelText(/Estimated weight of your spouse’s pro-gear/)).toBeNull();
      await userEvent.click(hasProGear);

      await waitFor(() => {
        expect(screen.getByLabelText(/Estimated weight of your pro-gear/)).toBeInstanceOf(HTMLInputElement);
        expect(screen.getByLabelText(/Estimated weight of your spouse’s pro-gear/)).toBeInstanceOf(HTMLInputElement);
      });
    });

    it('displays secondary gun safe weight inputs when hasGunSafe is true (ff on)', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

      await render(<EstimatedWeightsProGearForm {...defaultProps} />);
      const hasGunSafe = await screen.getByTestId('hasGunSafeYes');
      expect(screen.queryByLabelText(/Estimated weight of your gun safe/)).toBeNull();
      await userEvent.click(hasGunSafe);

      await waitFor(() => {
        expect(screen.getByLabelText(/Estimated weight of your gun safe/)).toBeInstanceOf(HTMLInputElement);
        expect(
          screen.queryByText(
            /The government authorizes the shipment of a gun safe up to 500 lbs. This is not charged against the authorized weight entitlement. The weight entitlement is charged for any weight over 500 lbs. The additional 500 lbs gun safe weight entitlement cannot be applied if a customer's overall entitlement is already at the 18,000 lbs maximum./,
          ),
        ).toBeInTheDocument();
      });
    });
  });

  describe('pull values from the ppm shipment when available', () => {
    it('renders blank form on load', async () => {
      render(<EstimatedWeightsProGearForm {...mtoShipmentProps} />);

      await waitFor(() => {
        expect(screen.getByLabelText(/Estimated weight of this PPM shipment/).value).toBe('4,000');
      });
      expect(screen.getByTestId('hasProGearYes').value).toBe('true');
      expect(screen.getByLabelText(/Estimated weight of your pro-gear/).value).toBe('1,000');
      expect(screen.getByLabelText(/Estimated weight of your spouse’s pro-gear/).value).toBe('100');
    });

    it('renders blank form on load with gun safe (ff on)', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      await render(<EstimatedWeightsProGearForm {...mtoShipmentPropsWithGunSafe} />);

      await waitFor(() => {
        expect(screen.getByLabelText(/Estimated weight of this PPM shipment/).value).toBe('4,000');
      });
      expect(screen.getByTestId('hasProGearYes').value).toBe('true');
      expect(screen.getByTestId('hasGunSafeYes').value).toBe('true');
      expect(screen.getByLabelText(/Estimated weight of your pro-gear/).value).toBe('1,000');
      expect(screen.getByLabelText(/Estimated weight of your spouse’s pro-gear/).value).toBe('100');
      expect(screen.getByLabelText(/Estimated weight of your gun safe/).value).toBe('500');
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

      const inputHasProGear = screen.getByTestId('hasProGearYes');

      await userEvent.click(inputHasProGear);

      const selfProGear = screen.getByLabelText(/Estimated weight of your pro-gear/);

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

    it('marks gun safe input as required when conditionally displayed (ff on)', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

      await render(<EstimatedWeightsProGearForm {...defaultProps} />);

      const inputHasGunSafe = screen.getByTestId('hasGunSafeYes');

      await userEvent.click(inputHasGunSafe);

      const gunSafeInput = screen.getByLabelText(/Estimated weight of your gun safe/);

      await userEvent.click(gunSafeInput);
      await userEvent.tab();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();

        const requiredAlerts = screen.getByRole('alert');

        expect(requiredAlerts).toHaveTextContent('Required');
      });
    });

    it('display error message for gun safe input above 500 lbs  (ff on)', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

      await render(<EstimatedWeightsProGearForm {...defaultProps} />);

      const inputHasGunSafe = screen.getByTestId('hasGunSafeYes');

      await userEvent.click(inputHasGunSafe);

      const gunSafeInput = screen.getByLabelText(/Estimated weight of your gun safe/);

      await userEvent.click(gunSafeInput);
      await userEvent.type(gunSafeInput, String(501));
      await userEvent.tab();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();

        const requiredAlerts = screen.getByRole('alert');

        expect(requiredAlerts).toHaveTextContent('Enter a weight 500 lbs or less');
      });
    });
  });
});

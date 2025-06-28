import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
// import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';
import { act } from 'react-dom/test-utils';

import ShipmentWeight from './ShipmentWeight';

import { isBooleanFlagEnabled } from 'utils/featureFlags';

jest.mock('../../../utils/featureFlags', () => ({
  ...jest.requireActual('../../../utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

describe('components/Office/ShipmentWeight', () => {
  it('defaults to customer not using Pro-gear or gun safe', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    await act(async () => {
      render(
        <Formik initialValues={{}}>
          <ShipmentWeight />
        </Formik>,
      );
    });

    expect(screen.getByTestId('hasProGearYes')).not.toBeChecked();
    expect(screen.getByTestId('hasProGearNo')).toBeChecked();

    expect(screen.queryByLabelText('Estimated pro-gear weight')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('Estimated spouse pro-gear weight')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('Estimated gun safe weight')).not.toBeInTheDocument();
    expect(
      screen.queryByText(
        `The government authorizes the shipment of a gun safe up to 500 lbs. The weight entitlement is charged for any weight over 500 lbs. The additional 500 lbs gun safe weight entitlement cannot be applied if a customer's overall entitlement is already at the 18,000 lbs maximum.`,
      ),
    ).not.toBeInTheDocument();
  });

  it('displays estimated weight and pro-gear data', async () => {
    render(
      <Formik
        initialValues={{
          hasProGear: true,
          estimatedWeight: '4000',
          proGearWeight: '3000',
          spouseProGearWeight: '2000',
        }}
      >
        <ShipmentWeight />
      </Formik>,
    );
    await waitFor(() => {
      expect(screen.getByTestId('hasProGearYes')).toBeChecked();
      expect(screen.getByTestId('hasProGearNo')).not.toBeChecked();

      expect(screen.queryByLabelText('Estimated pro-gear weight')).toBeInTheDocument();
      expect(screen.queryByLabelText('Estimated spouse pro-gear weight')).toBeInTheDocument();
    });
  });

  it('displays gun safe data', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    await act(async () => {
      render(
        <Formik
          initialValues={{
            hasGunSafe: true,
            estimatedWeight: '4000',
            gunSafeWeight: '455',
          }}
        >
          <ShipmentWeight />
        </Formik>,
      );
    });
    await waitFor(() => {
      expect(screen.getByTestId('hasGunSafeYes')).toBeInTheDocument();
      expect(screen.getByTestId('hasGunSafeYes')).toBeChecked();

      expect(screen.queryByLabelText('Estimated gun safe weight')).toBeInTheDocument();
      expect(
        screen.queryByText(
          `The government authorizes the shipment of a gun safe up to 500 lbs. The weight entitlement is charged for any weight over 500 lbs. The additional 500 lbs gun safe weight entitlement cannot be applied if a customer's overall entitlement is already at the 18,000 lbs maximum.`,
        ),
      ).toBeInTheDocument();
    });
  });
});

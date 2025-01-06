/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { screen } from '@testing-library/react';

import { isBooleanFlagEnabled } from '../../utils/featureFlags';
import { FEATURE_FLAG_KEYS } from '../../shared/constants';

import { MovingInfo } from './MovingInfo';

import { renderWithRouterProp } from 'testUtils';
import { customerRoutes } from 'constants/routes';

jest.mock('../../utils/featureFlags', () => ({
  ...jest.requireActual('../../utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve()),
}));

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectUbAllowance: jest.fn(),
}));

describe('MovingInfo component', () => {
  const testProps = {
    entitlementWeight: 7000,
    fetchLatestOrders: jest.fn(),
    serviceMemberId: '1234567890',
  };
  const routingOptions = {
    path: customerRoutes.SHIPMENT_MOVING_INFO_PATH,
    params: { moveId: 'testMove123' },
  };

  it('renders the expected content', () => {
    renderWithRouterProp(<MovingInfo {...testProps} />, routingOptions);

    expect(
      screen.getByRole('heading', { level: 1, name: 'Things to know about selecting shipments' }),
    ).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: /7,000 lbs/ })).toBeInTheDocument();
    expect(screen.getAllByRole('heading').length).toBe(6);
  });

  it('renders with no errors when entitlement weight is 0', () => {
    renderWithRouterProp(<MovingInfo {...testProps} entitlementWeight={0} />, routingOptions);

    expect(
      screen.getByRole('heading', { level: 1, name: 'Things to know about selecting shipments' }),
    ).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: /0 lbs/ })).toBeInTheDocument();
    expect(screen.getAllByRole('heading').length).toBe(6);
  });

  it('feature flag enable show headers for PPM', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

    const wrapper = renderWithRouterProp(<MovingInfo {...testProps} entitlementWeight={0} />, routingOptions);
    await wrapper;
    expect(
      screen.getByRole('heading', { name: /You still have the option to move some of your belongings yourself./ }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole('heading', { name: /You can get paid for any household goods you move yourself./ }),
    ).toBeInTheDocument();
    expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.PPM);
  });

  it('feature flag disable hides headers for PPM', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(false));

    const wrapper = renderWithRouterProp(<MovingInfo {...testProps} entitlementWeight={0} />, routingOptions);
    await wrapper;

    expect(screen.queryByText('You still have the option to move some of your belongings yourself.')).toBeNull();
    expect(screen.queryByText('You can get paid for any household goods you move yourself.')).toBeNull();

    expect(isBooleanFlagEnabled).toBeCalledWith(FEATURE_FLAG_KEYS.PPM);
  });

  it('when allowed UB allowance, they see the UB allowance section', async () => {
    const wrapper = renderWithRouterProp(
      <MovingInfo {...testProps} entitlementWeight={8000} ubAllowance={500} />,
      routingOptions,
    );
    await wrapper;

    expect(screen.queryByText('You can move up to 500 lbs of unaccompanied baggage in this move.')).toBeInTheDocument();
  });
});

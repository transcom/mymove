/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { screen } from '@testing-library/react';

import { MovingInfo } from './MovingInfo';

import { renderWithRouterProp } from 'testUtils';
import { customerRoutes } from 'constants/routes';

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
});

/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen } from '@testing-library/react';

import { MovingInfo } from './MovingInfo';

describe('MovingInfo component', () => {
  const testProps = {
    entitlementWeight: 7000,
    fetchLatestOrders: jest.fn(),
    serviceMemberId: '1234567890',
    match: {
      params: { moveId: 'testMove123' },
    },
    location: {},
    history: { push: jest.fn() },
  };

  it('renders the expected content', () => {
    render(<MovingInfo {...testProps} />);

    expect(
      screen.getByRole('heading', { level: 1, name: 'Things to know about selecting shipments' }),
    ).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: /7,000 lbs/ })).toBeInTheDocument();
    expect(screen.getAllByRole('heading').length).toBe(6);
  });

  it('renders with no errors when entitlement weight is 0', () => {
    render(<MovingInfo {...testProps} entitlementWeight={0} />);

    expect(
      screen.getByRole('heading', { level: 1, name: 'Things to know about selecting shipments' }),
    ).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: /0 lbs/ })).toBeInTheDocument();
    expect(screen.getAllByRole('heading').length).toBe(6);
  });
});

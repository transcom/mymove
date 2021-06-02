/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render } from '@testing-library/react';

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
    const { queryByText, queryAllByRole } = render(<MovingInfo {...testProps} />);

    expect(queryByText('Things to know about selecting shipments')).toBeInTheDocument();
    expect(queryByText(/7,000 lbs/)).toBeInTheDocument();
    expect(queryAllByRole('heading').length).toBe(6);
  });

  it('renders with no errors when entitlement weight is 0', () => {
    const { queryByText, queryAllByRole } = render(<MovingInfo {...testProps} entitlementWeight={0} />);

    expect(queryByText('Things to know about selecting shipments')).toBeInTheDocument();
    expect(queryAllByRole('heading').length).toBe(5);
  });
});

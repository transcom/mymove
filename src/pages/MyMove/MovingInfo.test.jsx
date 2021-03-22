/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render } from '@testing-library/react';

import { MovingInfo } from './MovingInfo';

describe('MovingInfo component', () => {
  const testProps = {
    entitlementWeight: 7000,
    fetchLatestOrders: jest.fn(),
    serviceMemberId: '1234567890',
  };

  it('renders the expected content', () => {
    const { queryAllByTestId, queryByText } = render(<MovingInfo {...testProps} />);

    expect(queryByText('MovingInfo')).toBeInTheDocument();
    expect(queryByText('Tips for planning your shipments')).toBeInTheDocument();
    expect(queryAllByTestId('shipmentsHeader')).toBeInTheDocument();
    expect(queryByText('7,000 lbs')).toBeInTheDocument();
    expect(queryAllByTestId('shipmentsSubHeader').length).toBe(4);
  });

  it('renders with no errors when entitlement weight is 0', () => {
    const { queryAllByTestId, queryByText } = render(<MovingInfo {...testProps} entitlementWeight={0} />);

    expect(queryByText('MovingInfo')).toBeInTheDocument();
    expect(queryByText('Tips for planning your shipments')).toBeInTheDocument();
    expect(queryAllByTestId('shipmentsHeader')).toBeInTheDocument();
    expect(queryAllByTestId('shipmentsAlert').length).toBe(0);
    expect(queryAllByTestId('shipmentsSubHeader').length).toBe(4);
  });
});

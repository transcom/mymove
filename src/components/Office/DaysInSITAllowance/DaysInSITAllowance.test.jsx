import React from 'react';
import { render, screen } from '@testing-library/react';

import DaysInSITAllowance from './DaysInSITAllowance';

describe('DaysInSITAllowance Component', () => {
  const shipmentPaymentSITBalance = {
    previouslyBilledDays: 30,
    previouslyBilledEndDate: '2021-06-08',
    pendingSITDaysInvoiced: 60,
    pendingBilledEndDate: '2021-08-08',
    pendingBilledStartDate: '2021-07-08',
    totalSITDaysAuthorized: 120,
    totalSITDaysRemaining: 30,
    totalSITEndDate: '2021-09-08',
  };

  const pendingShipmentSITBalance = {
    previouslyBilledDays: 0,
    pendingSITDaysInvoiced: 60,
    pendingBilledStartDate: '2021-07-08',
    pendingBilledEndDate: '2021-08-08',
    totalSITDaysAuthorized: 120,
    totalSITDaysRemaining: 30,
    totalSITEndDate: '2021-09-08',
  };

  it('renders the billed, pending, and total SIT balances', () => {
    render(<DaysInSITAllowance shipmentPaymentSITBalance={shipmentPaymentSITBalance} />);

    expect(screen.getByText('Prev. billed & accepted days')).toBeInTheDocument();
    // due to the fragments using getByText here doesn't work, another option would be create a function that renders a
    // single string fragment in the component
    expect(screen.getByTestId('previouslyBilled')).toHaveTextContent('30');

    expect(screen.getByText('Payment start - end date')).toBeInTheDocument();
    expect(screen.getByTestId('pendingBilledStartEndDate')).toHaveTextContent('08 Jul 2021 - 08 Aug 2021');

    expect(screen.getByText('Total days of SIT approved')).toBeInTheDocument();
    expect(screen.getByText('120')).toBeInTheDocument();

    expect(screen.getByText('Total approved days remaining')).toBeInTheDocument();
    expect(screen.getByTestId('totalRemaining')).toHaveTextContent('30');
  });

  it('renders zero when no SIT days were previously billed', () => {
    render(<DaysInSITAllowance shipmentPaymentSITBalance={pendingShipmentSITBalance} />);

    expect(screen.getByText('Prev. billed & accepted days')).toBeInTheDocument();
    expect(screen.getByTestId('previouslyBilled')).toHaveTextContent('0');

    expect(screen.getByText('Payment start - end date')).toBeInTheDocument();
    expect(screen.getByTestId('pendingBilledStartEndDate')).toHaveTextContent('08 Jul 2021 - 08 Aug 2021');

    expect(screen.getByText('Total days of SIT approved')).toBeInTheDocument();
    expect(screen.getByText('120')).toBeInTheDocument();

    expect(screen.getByText('Total approved days remaining')).toBeInTheDocument();
    expect(screen.getByTestId('totalRemaining')).toHaveTextContent('30');
  });
});

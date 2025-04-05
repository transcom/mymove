import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';

import PaymentDetails from './PaymentDetails';

describe('PaymentDetails', () => {
  describe('for each changed value', () => {
    const context = [
      {
        name: 'Test Service',
        price: '10123',
        status: 'APPROVED',
      },
      {
        name: 'Domestic uncrating',
        price: '5555',
        status: 'APPROVED',
      },
    ];
    it.each([
      ['Test Service', '101.23'],
      ['Domestic uncrating', '55.55'],
    ])('it renders %s: %s', (displayName, value) => {
      render(<PaymentDetails context={context} />);

      expect(screen.getByText(displayName)).toBeInTheDocument();

      expect(screen.getByText(value, { exact: false })).toBeInTheDocument();
    });
  });

  it('Sums the values', async () => {
    const context = [
      {
        name: 'Test Service',
        price: '10123',
        status: 'APPROVED',
      },
      {
        name: 'Domestic uncrating',
        price: '5555',
        status: 'APPROVED',
      },
    ];

    render(<PaymentDetails context={context} />);

    expect(screen.getByText(156.78, { exact: false })).toBeInTheDocument();
  });

  describe('rejected service items', () => {
    const context = [
      {
        name: 'Domestic uncrating',
        price: '5555',
        status: 'DENIED',
        rejection_reason: 'some reason',
      },
    ];
    it('renders a rejected service item and its rejection reason', async () => {
      render(<PaymentDetails context={context} />);

      expect(screen.getByText('Domestic uncrating')).toBeInTheDocument();

      expect(screen.getByText('Rejection Reason:')).toBeInTheDocument();
      await waitFor(() => {
        screen.getByText('Rejection Reason:').click();
      });
      expect(screen.getByText('some reason')).toBeVisible();
    });
  });
});

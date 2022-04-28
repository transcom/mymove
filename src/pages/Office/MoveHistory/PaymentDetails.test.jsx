import React from 'react';
import { render, screen } from '@testing-library/react';

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
});

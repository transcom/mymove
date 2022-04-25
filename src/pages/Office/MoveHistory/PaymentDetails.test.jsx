import React from 'react';
import { render, screen } from '@testing-library/react';

import PaymentDetails from './PaymentDetails';

describe('PaymentDetails', () => {
  describe('for each changed value', () => {
    const context = [
      {
        name: 'Test Service',
        price: '101',
        status: 'APPROVED',
      },
      {
        name: 'Domestic uncrating',
        price: '55',
        status: 'APPROVED',
      },
    ];
    it.each([
      ['Test Service', '101'],
      ['Domestic uncrating', '55'],
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
        price: '101',
        status: 'APPROVED',
      },
      {
        name: 'Domestic uncrating',
        price: '55',
        status: 'APPROVED',
      },
    ];

    render(<PaymentDetails context={context} />);

    expect(screen.getByText(156, { exact: false })).toBeInTheDocument();
  });
});

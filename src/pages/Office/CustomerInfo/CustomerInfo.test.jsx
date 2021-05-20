import React from 'react';
import { render, screen } from '@testing-library/react';

import CustomerInfo from './CustomerInfo';

import { MockProviders } from 'testUtils';

const mockCustomer = {
  backup_contact: {
    email: 'backup@mail.com',
    name: 'Jane Backup',
    phone: '555-555-1234',
  },
  current_address: {
    city: 'Beverly Hills',
    country: 'US',
    postal_code: '90210',
    state: 'CA',
    street_address_1: '123 Any Street',
  },
  email: 'john_doe@mail.com',
  first_name: 'John',
  last_name: 'Doe',
  middle_name: 'Quincey',
  suffix: 'Jr.',
  phone: '123-444-3434',
};

describe('CustomerInfo', () => {
  it('populates initial field values', () => {
    render(
      <MockProviders initialEntries={['moves/CDG3TR/customer-info']}>
        <CustomerInfo customer={mockCustomer} isLoading={false} isError={false} />{' '}
      </MockProviders>,
    );
    expect(screen.getByLabelText('First name').value).toEqual('John');
  });
});

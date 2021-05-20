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
    expect(screen.getByLabelText('First name').value).toEqual(mockCustomer.first_name);
    expect(screen.getByLabelText(/Middle name/i).value).toEqual(mockCustomer.middle_name);
    expect(screen.getByLabelText('Last name').value).toEqual(mockCustomer.last_name);
    expect(screen.getByLabelText(/Suffix/i).value).toEqual(mockCustomer.suffix);
    // to get around the two inputs labeled "Phone" on the screen
    expect(screen.getByDisplayValue(mockCustomer.phone).value).toEqual(mockCustomer.phone);
    expect(screen.getByDisplayValue(mockCustomer.backup_contact.phone).value).toEqual(
      mockCustomer.backup_contact.phone,
    );
    // to get around the two inputs labeled "Email" on the screen
    expect(screen.getByDisplayValue(mockCustomer.email).value).toEqual(mockCustomer.email);
    expect(screen.getByDisplayValue(mockCustomer.backup_contact.email).value).toEqual(
      mockCustomer.backup_contact.email,
    );
    expect(screen.getByLabelText('Address 1').value).toEqual(mockCustomer.current_address.street_address_1);
    expect(screen.getByLabelText('City').value).toEqual(mockCustomer.current_address.city);
    expect(screen.getByLabelText('State').value).toEqual(mockCustomer.current_address.state);
    expect(screen.getByLabelText('ZIP').value).toEqual(mockCustomer.current_address.postal_code);
    expect(screen.getByLabelText('Name').value).toEqual(mockCustomer.backup_contact.name);
  });
});

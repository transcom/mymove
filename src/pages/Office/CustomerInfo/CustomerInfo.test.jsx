import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { generatePath } from 'react-router';
import userEvent from '@testing-library/user-event';

import CustomerInfo from './CustomerInfo';

import { MockProviders } from 'testUtils';
import { updateCustomerInfo } from 'services/ghcApi';
import { servicesCounselingRoutes } from 'constants/routes';

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  updateCustomerInfo: jest.fn(),
}));

const mockRequestedMoveCode = 'LR4T8V';
const customerInfoEditURL = generatePath(servicesCounselingRoutes.CUSTOMER_INFO_EDIT_PATH, {
  moveCode: mockRequestedMoveCode,
});

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn().mockReturnValue({ moveCode: 'LR4T8V' }),
}));

const mockCustomer = {
  backup_contact: {
    email: 'backup@mail.com',
    name: 'Jane Backup',
    phone: '555-555-1234',
  },
  current_address: {
    city: 'Beverly Hills',
    country: 'US',
    postalCode: '90210',
    state: 'CA',
    streetAddress1: '123 Any Street',
  },
  email: 'john_doe@mail.com',
  first_name: 'John',
  last_name: 'Doe',
  middle_name: 'Quincey',
  suffix: 'Jr.',
  phone: '123-444-3434',
};

const loadingReturnValue = {
  isLoading: true,
  isError: false,
  isSuccess: false,
};

const errorReturnValue = {
  isLoading: false,
  isError: true,
  isSuccess: false,
};

let mockUpdate;

describe('CustomerInfo', () => {
  beforeEach(() => {
    mockUpdate = jest.fn();
  });

  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      updateCustomerInfo.mockReturnValue(loadingReturnValue);

      render(
        <MockProviders initialEntries={[customerInfoEditURL]}>
          <CustomerInfo customer={mockCustomer} onUpdate={mockUpdate} ordersId="abc123" isLoading isError={false} />{' '}
        </MockProviders>,
      );

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      updateCustomerInfo.mockReturnValue(errorReturnValue);

      render(
        <MockProviders initialEntries={[customerInfoEditURL]}>
          <CustomerInfo customer={mockCustomer} onUpdate={mockUpdate} ordersId="abc123" isLoading={false} isError />{' '}
        </MockProviders>,
      );

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  it('populates initial field values', async () => {
    render(
      <MockProviders initialEntries={[customerInfoEditURL]}>
        <CustomerInfo
          customer={mockCustomer}
          onUpdate={mockUpdate}
          ordersId="abc123"
          isLoading={false}
          isError={false}
        />
      </MockProviders>,
    );
    await waitFor(() => {
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
      expect(screen.getByLabelText('Address 1').value).toEqual(mockCustomer.current_address.streetAddress1);
      expect(screen.getByLabelText('City').value).toEqual(mockCustomer.current_address.city);
      expect(screen.getByLabelText('State').value).toEqual(mockCustomer.current_address.state);
      expect(screen.getByLabelText('ZIP').value).toEqual(mockCustomer.current_address.postalCode);
      expect(screen.getByLabelText('Name').value).toEqual(mockCustomer.backup_contact.name);
    });
  });

  it('calls onUpdate prop with success on successful form submission', async () => {
    updateCustomerInfo.mockImplementation(() => Promise.resolve({ customer: { customerId: '123' } }));
    render(
      <MockProviders initialEntries={[customerInfoEditURL]}>
        <CustomerInfo
          customer={mockCustomer}
          onUpdate={mockUpdate}
          ordersId="abc123"
          isLoading={false}
          isError={false}
        />
      </MockProviders>,
    );
    const saveBtn = screen.getByRole('button', { name: 'Save' });
    await userEvent.click(saveBtn);

    await waitFor(() => {
      expect(mockUpdate).toHaveBeenCalledWith('success');
    });
  });

  it('calls onUpdate prop with error on unsuccessful form submission', async () => {
    jest.spyOn(console, 'error').mockImplementation(() => {});
    updateCustomerInfo.mockImplementation(() => Promise.reject());

    render(
      <MockProviders initialEntries={[customerInfoEditURL]}>
        <CustomerInfo
          customer={mockCustomer}
          ordersId="abc123"
          isLoading={false}
          isError={false}
          onUpdate={mockUpdate}
        />
      </MockProviders>,
    );

    const saveBtn = screen.getByRole('button', { name: 'Save' });
    await userEvent.click(saveBtn);

    await waitFor(async () => {
      await expect(mockUpdate).toHaveBeenCalledWith('error');
    });
  });
});

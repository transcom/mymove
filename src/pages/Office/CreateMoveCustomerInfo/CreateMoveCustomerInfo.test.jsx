import React from 'react';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import CreateMoveCustomerInfo from './CreateMoveCustomerInfo';

import { renderWithProviders } from 'testUtils';
import { updateCustomerInfo } from 'services/ghcApi';
import { servicesCounselingRoutes } from 'constants/routes';
import { useCustomerQuery } from 'hooks/queries';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));
const routingParams = { customerId: '8604447b-cbfc-4d59-a9a1-dec219eb2046' };
const mockRoutingConfig = {
  path: servicesCounselingRoutes.BASE_CUSTOMERS_CUSTOMER_INFO_PATH,
  params: routingParams,
};

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  updateCustomerInfo: jest.fn(),
}));

jest.mock('hooks/queries', () => ({
  useCustomerQuery: jest.fn(),
}));

const useCustomerQueryReturnValue = {
  customerData: {
    backup_contact: {
      email: 'backup@mail.com',
      name: 'Jane Backup',
      phone: '555-555-1234',
    },
    backupAddress: {
      city: 'Great Falls',
      country: 'US',
      postalCode: '59402',
      state: 'MT',
      streetAddress1: '446 South Ave',
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
    phone: '223-444-3434',
    cacValidated: true,
  },
};

const loadingReturnValue = {
  ...useCustomerQueryReturnValue,
  isLoading: true,
  isError: false,
  isSuccess: false,
};

const errorReturnValue = {
  ...useCustomerQueryReturnValue,
  isLoading: false,
  isError: true,
  isSuccess: false,
};

describe('CreateMoveCustomerInfo', () => {
  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useCustomerQuery.mockReturnValue(loadingReturnValue);

      renderWithProviders(<CreateMoveCustomerInfo />, mockRoutingConfig);

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useCustomerQuery.mockReturnValue(errorReturnValue);

      renderWithProviders(<CreateMoveCustomerInfo />, mockRoutingConfig);

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  it('populates initial field values', async () => {
    useCustomerQuery.mockReturnValue(useCustomerQueryReturnValue);

    renderWithProviders(<CreateMoveCustomerInfo />, mockRoutingConfig);
    const { customerData } = useCustomerQueryReturnValue;
    await waitFor(() => {
      expect(screen.getByLabelText('First name').value).toEqual(customerData.first_name);
      expect(screen.getByLabelText(/Middle name/i).value).toEqual(customerData.middle_name);
      expect(screen.getByLabelText('Last name').value).toEqual(customerData.last_name);
      expect(screen.getByLabelText(/Suffix/i).value).toEqual(customerData.suffix);
      // to get around the two inputs labeled "Phone" on the screen
      expect(screen.getByDisplayValue(customerData.phone).value).toEqual(customerData.phone);
      expect(screen.getByDisplayValue(customerData.backup_contact.phone).value).toEqual(
        customerData.backup_contact.phone,
      );
      // to get around the two inputs labeled "Email" on the screen
      expect(screen.getByDisplayValue(customerData.email).value).toEqual(customerData.email);
      expect(screen.getByDisplayValue(customerData.backup_contact.email).value).toEqual(
        customerData.backup_contact.email,
      );
      expect(screen.getByDisplayValue('123 Any Street').value).toEqual(customerData.current_address.streetAddress1);
      expect(screen.getByText('Beverly Hills')).toHaveTextContent(customerData.current_address.city);
      expect(screen.getByText('CA')).toHaveTextContent(customerData.current_address.state);
      expect(screen.getByText('90210')).toHaveTextContent(customerData.current_address.postalCode);
      expect(screen.getByDisplayValue('Jane Backup').value).toEqual(customerData.backup_contact.name);
    });
  });

  it('calls updateCustomerInfo on submission', async () => {
    useCustomerQuery.mockReturnValue(useCustomerQueryReturnValue);
    updateCustomerInfo.mockImplementation(() => Promise.resolve({ customer: { customerId: '123' } }));
    renderWithProviders(<CreateMoveCustomerInfo />, mockRoutingConfig);

    const saveBtn = screen.getByRole('button', { name: 'Save' });
    await userEvent.click(saveBtn);

    await waitFor(() => {
      expect(updateCustomerInfo).toHaveBeenCalled();
    });
  });
});

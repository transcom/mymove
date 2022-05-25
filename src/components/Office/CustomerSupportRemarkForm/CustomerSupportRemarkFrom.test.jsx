import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import CustomerSupportRemarkForm from './CustomerSupportRemarkForm';

import { MockProviders } from 'testUtils';
import * as api from 'services/ghcApi';
import { roleTypes } from 'constants/userRoles';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn().mockReturnValue({ moveCode: 'LR4T8V' }),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const testState = {
  auth: {
    activeRole: roleTypes.QAE_CSR,
    isLoading: false,
    isLoggedIn: true,
  },
  entities: {
    user: {
      userId123: {
        id: 'userId123',
        roles: [{ roleType: roleTypes.QAE_CSR }],
        office_user: {
          first_name: 'Amanda',
          last_name: 'Gorman',
          transportation_office: {
            gbloc: 'ABCD',
          },
        },
      },
    },
  },
};

describe('CustomerSupportRemarkForm', () => {
  it('renders the form', async () => {
    render(
      <MockProviders initialState={testState}>
        <CustomerSupportRemarkForm />
      </MockProviders>,
    );

    expect(await screen.findByTestId('form')).toBeInTheDocument();
    expect(await screen.findByTestId('textarea')).toBeInTheDocument();
    expect(await screen.findByTestId('button')).toBeInTheDocument();
  });

  it('submits the form with expected data', async () => {
    // Spy on and mock mutation function
    const mutationSpy = jest.spyOn(api, 'createCustomerSupportRemarkForMove').mockImplementation(() => {});
    render(
      <MockProviders initialState={testState}>
        <CustomerSupportRemarkForm />
      </MockProviders>,
    );

    // Type in the textarea
    userEvent.type(await screen.findByTestId('textarea'), 'Test Remark');
    await waitFor(() => {
      expect(screen.getByTestId('button').hasAttribute('disabled')).toBeFalsy();
    });

    // Submit the form
    await waitFor(() => {
      fireEvent.click(screen.getByRole('button', { name: 'Save' }));
    });

    // Ensure the expected mutation was called with expected data
    expect(mutationSpy).toHaveBeenCalledTimes(1);
    expect(mutationSpy).toHaveBeenCalledWith({
      locator: 'LR4T8V',
      content: 'Test Remark',
    });
  });

  it('will not submit empty remarks', async () => {
    // Spy on and mock mutation function
    const mutationSpy = jest.spyOn(api, 'createCustomerSupportRemarkForMove').mockImplementation(() => {});

    render(
      <MockProviders initialState={testState}>
        <CustomerSupportRemarkForm />
      </MockProviders>,
    );

    // Submit the empty form
    await waitFor(() => {
      fireEvent.click(screen.getByRole('button', { name: 'Save' }));
    });

    // Ensure the expected mutation was called with expected data
    expect(mutationSpy).toHaveBeenCalledTimes(0);
  });
});

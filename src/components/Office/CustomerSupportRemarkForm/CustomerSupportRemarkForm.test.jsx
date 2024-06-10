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

const qaeTestState = {
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

const csrTestState = {
  auth: {
    activeRole: roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE,
    isLoading: false,
    isLoggedIn: true,
  },
  entities: {
    user: {
      userId123: {
        id: 'userId123',
        roles: [{ roleType: roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE }],
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
  it('qae can render the form', async () => {
    render(
      <MockProviders initialState={qaeTestState}>
        <CustomerSupportRemarkForm />
      </MockProviders>,
    );

    expect(await screen.findByTestId('form')).toBeInTheDocument();
    expect(await screen.findByTestId('textarea')).toBeInTheDocument();
    expect(await screen.findByTestId('button')).toBeInTheDocument();
  });
  it('csr can render the form', async () => {
    render(
      <MockProviders initialState={csrTestState}>
        <CustomerSupportRemarkForm />
      </MockProviders>,
    );

    expect(await screen.findByTestId('form')).toBeInTheDocument();
    expect(await screen.findByTestId('textarea')).toBeInTheDocument();
    expect(await screen.findByTestId('button')).toBeInTheDocument();
  });

  it('qae can submit the form with expected data', async () => {
    // Spy on and mock mutation function
    const mutationSpy = jest.spyOn(api, 'createCustomerSupportRemarkForMove').mockImplementation(() => {});
    render(
      <MockProviders initialState={qaeTestState}>
        <CustomerSupportRemarkForm />
      </MockProviders>,
    );

    // Type in the textarea
    await userEvent.type(await screen.findByTestId('textarea'), 'Test Remark');
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
      body: {
        content: 'Test Remark',
      },
      locator: 'LR4T8V',
    });
  });

  it('csr can submit the form with expected data', async () => {
    // Spy on and mock mutation function
    const mutationSpy = jest.spyOn(api, 'createCustomerSupportRemarkForMove').mockImplementation(() => {});
    render(
      <MockProviders initialState={csrTestState}>
        <CustomerSupportRemarkForm />
      </MockProviders>,
    );

    // Type in the textarea
    await userEvent.type(await screen.findByTestId('textarea'), 'Test Remark');
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
      body: {
        content: 'Test Remark',
      },
      locator: 'LR4T8V',
    });
  });

  it('qae will not submit empty remarks', async () => {
    // Spy on and mock mutation function
    const mutationSpy = jest.spyOn(api, 'createCustomerSupportRemarkForMove').mockImplementation(() => {});

    render(
      <MockProviders initialState={qaeTestState}>
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

  it('csr will not submit empty remarks', async () => {
    // Spy on and mock mutation function
    const mutationSpy = jest.spyOn(api, 'createCustomerSupportRemarkForMove').mockImplementation(() => {});

    render(
      <MockProviders initialState={csrTestState}>
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

  it('disables the submit button when the move is locked', async () => {
    // Spy on and mock mutation function
    const isMoveLocked = true;
    render(
      <MockProviders initialState={qaeTestState}>
        <CustomerSupportRemarkForm isMoveLocked={isMoveLocked} />
      </MockProviders>,
    );

    // Type in the textarea, button should still be disabled
    await userEvent.type(await screen.findByTestId('textarea'), 'Test Remark');
    expect(screen.getByTestId('button').hasAttribute('disabled')).toBeTruthy();
  });
});

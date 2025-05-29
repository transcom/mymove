import React from 'react';
import { screen, render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { act } from 'react-dom/test-utils';

import { ConnectedSelectApplication, roleLookupValues } from './MultiRoleSelectApplication';

import { MockProviders } from 'testUtils';
import { setActiveRole } from 'store/auth/actions';
import { configureStore } from 'shared/store';

jest.mock('store/auth/actions', () => ({
  ...jest.requireActual('store/auth/actions'),
  setActiveRole: jest.fn().mockImplementation(() => ({ type: '' })),
}));

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => jest.fn(),
}));

describe('MultiRoleSelectApplication component', () => {
  it('renders "no role" when there are no roles', async () => {
    const mockState = {
      entities: {
        user: [
          {
            inactiveRoles: [],
          },
        ],
      },
      auth: {
        activeRole: null,
      },
    };

    const mockStore = configureStore({});
    render(
      <MockProviders store={mockStore} initialState={mockState}>
        <ConnectedSelectApplication />
      </MockProviders>,
    );

    const roleLabel = await screen.findByText('Role:');
    expect(roleLabel).toBeInTheDocument();
  });

  it('renders the active role if one exists', async () => {
    const testUserRoles = [roleLookupValues.services_counselor];

    const mockState = {
      entities: {
        user: [
          {
            inactiveRoles: testUserRoles,
          },
        ],
      },
      auth: {
        activeRole: roleLookupValues.services_counselor.roleType,
      },
    };

    const mockStore = configureStore({});
    render(
      <MockProviders store={mockStore} initialState={mockState}>
        <ConnectedSelectApplication />
      </MockProviders>,
    );

    const roleLabel = await screen.findByText('Role:');
    expect(roleLabel).toBeInTheDocument();
  });

  it('renders options for each of the user’s roles, and does not render options for roles the user doesn’t have', async () => {
    const testUserRoles = [
      roleLookupValues.task_ordering_officer,
      roleLookupValues.task_invoicing_officer,
      roleLookupValues.services_counselor,
      roleLookupValues.qae,
      roleLookupValues.customer_service_representative,
    ];

    const mockState = {
      entities: {
        user: [
          {
            inactiveRoles: testUserRoles,
          },
        ],
      },
      auth: {
        activeRole: roleLookupValues.services_counselor.roleType,
      },
    };

    const mockStore = configureStore({});
    render(
      <MockProviders store={mockStore} initialState={mockState}>
        <ConnectedSelectApplication />
      </MockProviders>,
    );
    await Promise.all(
      testUserRoles.map(async ({ name }) => {
        const locatedOption = await screen.findByText(name);
        expect(locatedOption).toBeInTheDocument();
      }),
    );
  });

  it('handles setActiveRole with the selected role', async () => {
    const testUserRoles = [
      roleLookupValues.task_ordering_officer,
      roleLookupValues.task_invoicing_officer,
      roleLookupValues.qae,
      roleLookupValues.customer_service_representative,
    ];

    const mockState = {
      entities: {
        user: [
          {
            inactiveRoles: testUserRoles,
          },
        ],
      },
      auth: {
        activeRole: roleLookupValues.services_counselor.roleType,
      },
    };

    const mockStore = configureStore({});
    render(
      <MockProviders store={mockStore} initialState={mockState}>
        <ConnectedSelectApplication />
      </MockProviders>,
    );

    const user = userEvent.setup();

    waitFor(() => {
      expect(setActiveRole).toHaveBeenCalledWith(roleLookupValues.services_counselor.roleType);
    });

    await act(async () => {
      await Promise.all(
        testUserRoles.map(async ({ name }) => {
          const locatedOption = await screen.findByText(name);
          expect(locatedOption).toBeInTheDocument();
        }),
      );

      const dropdown = await screen.findByRole('combobox');
      const optionToSelect = roleLookupValues.task_invoicing_officer.roleType;
      await user.selectOptions(dropdown, optionToSelect);
      await user.tab();

      waitFor(() => {
        expect(setActiveRole).toHaveBeenCalledWith(optionToSelect);
      });
    });

    await act(async () => {
      const optionToCompare = roleLookupValues.task_invoicing_officer.name;
      const locatedOption = await screen.findByRole('option', { name: optionToCompare });
      expect(locatedOption.selected).toBe(true);
    });
  });
});

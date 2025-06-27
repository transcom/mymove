import React from 'react';
import { screen, render, waitFor, within } from '@testing-library/react';
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

    const roleLabel = await screen.findByText('none');
    expect(roleLabel).toBeInTheDocument();
  });

  it('renders the active role if one exists', async () => {
    const mockState = {
      entities: {
        user: [],
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

    expect(await screen.findByText('Role:')).toBeInTheDocument();
    expect(await screen.findByText(roleLookupValues.services_counselor.abbv)).toBeInTheDocument();
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
      testUserRoles.map(async ({ abbv }) => {
        const locatedOption = await screen.findByText(abbv);
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

    await waitFor(() => {
      expect(setActiveRole).toHaveBeenCalledWith(roleLookupValues.services_counselor.roleType);
    });

    const dropdown = await screen.findByRole('combobox');
    const optionToSelect = roleLookupValues.task_invoicing_officer.roleType;

    await act(async () => {
      await Promise.all(
        testUserRoles.map(async ({ abbv }) => {
          expect(await within(dropdown).findByText(abbv)).toBeInTheDocument();
        }),
      );

      await user.selectOptions(dropdown, optionToSelect);
    });

    await waitFor(() => {
      expect(setActiveRole).toHaveBeenLastCalledWith(optionToSelect);
    });

    const optionToCompare = roleLookupValues.task_invoicing_officer.roleType;
    expect(dropdown).toHaveValue(optionToCompare);
  });
});

import React from 'react';
import { screen, render, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { act } from 'react-dom/test-utils';

import { ConnectedSelectApplication, roleLookupValues } from './MultiRoleSelectApplication';

import { MockProviders } from 'testUtils';
import { setActiveRole } from 'store/auth/actions';
import { configureStore } from 'shared/store';
import { MULTI_SELECT_DROPDOWN_ARIA_TEXT } from 'utils/formatters';

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

  const roleDropdownScenarios = {
    rolesProvidedInProperOrder: {
      input: [
        roleLookupValues.services_counselor,
        roleLookupValues.task_ordering_officer,
        roleLookupValues.task_invoicing_officer,
        roleLookupValues.qae,
        roleLookupValues.customer_service_representative,
        roleLookupValues.gsr,
        roleLookupValues.headquarters,
        roleLookupValues.contracting_officer,
      ],
      expected: [
        roleLookupValues.services_counselor,
        roleLookupValues.task_ordering_officer,
        roleLookupValues.task_invoicing_officer,
        roleLookupValues.qae,
        roleLookupValues.customer_service_representative,
        roleLookupValues.gsr,
        roleLookupValues.headquarters,
        roleLookupValues.contracting_officer,
      ],
    },
    fourRoles: {
      input: [
        roleLookupValues.services_counselor,
        roleLookupValues.qae,
        roleLookupValues.customer_service_representative,
        roleLookupValues.gsr,
      ],
      expected: [
        roleLookupValues.services_counselor,
        roleLookupValues.qae,
        roleLookupValues.customer_service_representative,
        roleLookupValues.gsr,
      ],
    },
    rolesNotInTheCorrectOrder: {
      input: [
        roleLookupValues.task_ordering_officer,
        roleLookupValues.contracting_officer,
        roleLookupValues.headquarters,
        roleLookupValues.gsr,
        roleLookupValues.customer_service_representative,
        roleLookupValues.qae,
        roleLookupValues.task_invoicing_officer,
        roleLookupValues.services_counselor,
      ],
      expected: [
        roleLookupValues.services_counselor,
        roleLookupValues.task_ordering_officer,
        roleLookupValues.task_invoicing_officer,
        roleLookupValues.qae,
        roleLookupValues.customer_service_representative,
        roleLookupValues.gsr,
        roleLookupValues.headquarters,
        roleLookupValues.contracting_officer,
      ],
    },
    twoRoles: {
      input: [roleLookupValues.task_invoicing_officer, roleLookupValues.services_counselor],
      expected: [roleLookupValues.services_counselor, roleLookupValues.task_invoicing_officer],
    },
  };

  const roleLabelScenarios = {
    oneRoleAsServicesCounselor: {
      input: [roleLookupValues.services_counselor],
      expected: [roleLookupValues.services_counselor],
    },
    oneRoleAsHeadQuarters: {
      input: [roleLookupValues.headquarters],
      expected: [roleLookupValues.headquarters],
    },
  };

  const getOptionValues = (option) => {
    const optionValue = {
      roleType: option.value,
      abbv: option.textContent,
      name: option.getAttribute('aria-label'),
    };
    return optionValue;
  };

  it.each([roleLabelScenarios.oneRoleAsServicesCounselor, roleLabelScenarios.oneRoleAsHeadQuarters])(
    `properly displays the user role`,
    async ({ input, expected }) => {
      const expecting = expected[0];
      const [firstRole, ...otherRoles] = input;
      const mockState = {
        entities: {
          user: [
            {
              inactiveRoles: otherRoles,
            },
          ],
        },
        auth: {
          activeRole: firstRole.roleType,
        },
      };

      const mockStore = configureStore({});
      render(
        <MockProviders store={mockStore} initialState={mockState}>
          <ConnectedSelectApplication />
        </MockProviders>,
      );
      act(async () => {
        const labelElement = await screen.findByLabelText(MULTI_SELECT_DROPDOWN_ARIA_TEXT.label(firstRole.name));
        expect(labelElement).toBeInTheDocument();
        const linkValue = await screen.findByText(firstRole.abbv);
        expect(linkValue).toHaveTextContent(expecting.abbv);
      });
    },
  );

  it.each([
    roleDropdownScenarios.rolesProvidedInProperOrder,
    roleDropdownScenarios.rolesNotInTheCorrectOrder,
    roleDropdownScenarios.fourRoles,
    roleDropdownScenarios.twoRoles,
  ])('properly displays the order of user roles', async ({ input, expected }) => {
    const [firstRole] = input;
    const mockState = {
      entities: {
        user: [
          {
            inactiveRoles: input,
          },
        ],
      },
      auth: {
        activeRole: firstRole.roleType,
      },
    };

    const mockStore = configureStore({});
    render(
      <MockProviders store={mockStore} initialState={mockState}>
        <ConnectedSelectApplication />
      </MockProviders>,
    );

    const labelElement = await screen.findByLabelText(MULTI_SELECT_DROPDOWN_ARIA_TEXT.combobox);
    expect(labelElement).toBeInTheDocument();
    const selectElement = await screen.findByLabelText('User roles');
    const options = await within(selectElement).findAllByRole('option');
    const optionValues = options.map(getOptionValues);

    expect(expected).toEqual(optionValues);
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

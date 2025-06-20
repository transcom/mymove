import React from 'react';

import MultiRoleSelectApplication, { ConnectedSelectApplication, roleLookupValues } from './MultiRoleSelectApplication';
import style from './MultiRoleStoryDecorator.module.scss';

import { configureStore } from 'shared/store';
import { adminOfficeRoles } from 'constants/userRoles';
import { MockProviders } from 'testUtils';

const mockStore = configureStore({});

const roleMap = Object.fromEntries(adminOfficeRoles.map(({ roleType, name }) => [roleType, { roleType, name }]));

const optionInactiveRoles = Object.keys(roleMap);

export default {
  title: 'Office Components/MultiRoleSelect',
  decorators: [
    (Story) => (
      <div className={style.wrapper}>
        <Story />
      </div>
    ),
  ],
  component: MultiRoleSelectApplication,
  args: {},
  argTypes: {
    inactiveRoles: {
      options: optionInactiveRoles,
      mapping: roleMap,
      control: {
        type: 'multi-select',
        label: Object.fromEntries(Object.values(roleMap).map(({ roleType }) => [roleType, roleType])),
      },
    },
    activeRole: {
      control: { type: 'radio' },
      options: optionInactiveRoles,
    },
  },
};

export const MultiRoleUser = {
  render: ({
    activeRole = roleLookupValues.services_counselor,
    inactiveRoles = [
      roleLookupValues.services_counselor,
      roleLookupValues.task_ordering_officer,
      roleLookupValues.task_invoicing_officer,
      roleLookupValues.qae,
    ],
  }) => {
    const roles = inactiveRoles?.filter(({ roleType }) => roleType !== activeRole);
    const mockState = {
      auth: {
        activeRole,
      },
      entities: {
        user: [
          {
            inactiveRoles: roles,
          },
        ],
      },
    };
    return (
      <MockProviders store={mockStore} initialState={mockState}>
        <ConnectedSelectApplication />
      </MockProviders>
    );
  },
};

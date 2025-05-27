import React from 'react';

import MultiRoleSelectApplication from './MultiRoleSelectApplication';
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
  args: {
    setActiveRole: () => {},
  },
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
  render: (args) => (
    <MockProviders store={mockStore.store}>
      <MultiRoleSelectApplication {...args} />
    </MockProviders>
  ),
};

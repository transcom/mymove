import React from 'react';
import { mount } from 'enzyme';

import SelectApplication from './SelectApplication';

import { roleTypes } from 'constants/userRoles';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: jest.fn(),
  }),
}));

describe('SelectApplication component', () => {
  it('renders the active role if one exists', () => {
    const mockSetActiveRole = jest.fn();
    const wrapper = mount(<SelectApplication activeRole="myRole" userRoles={[]} setActiveRole={mockSetActiveRole} />);
    expect(wrapper.containsMatchingElement(<h2>Current role: myRole</h2>)).toEqual(true);
  });

  it('renders the first user role if there is no active role', () => {
    const mockSetActiveRole = jest.fn();
    const wrapper = mount(
      <SelectApplication
        userRoles={[{ roleType: 'myFirstRole' }, { roleType: 'myOtherRole' }]}
        setActiveRole={mockSetActiveRole}
      />,
    );
    expect(wrapper.containsMatchingElement(<h2>Current role: myFirstRole</h2>)).toEqual(true);
  });

  it('renders buttons for each of the user’s roles, and does not render buttons for roles the user doesn’t have', () => {
    const mockSetActiveRole = jest.fn();
    const wrapper = mount(
      <SelectApplication
        userRoles={[
          { roleType: roleTypes.TOO },
          { roleType: roleTypes.TIO },
          { roleType: roleTypes.SERVICES_COUNSELOR },
          { roleType: roleTypes.QAE_CSR },
        ]}
        setActiveRole={mockSetActiveRole}
      />,
    );

    expect(wrapper.containsMatchingElement(<button type="button">Select {roleTypes.TOO}</button>)).toEqual(true);
    expect(wrapper.containsMatchingElement(<button type="button">Select {roleTypes.TIO}</button>)).toEqual(true);
    expect(
      wrapper.containsMatchingElement(<button type="button">Select {roleTypes.SERVICES_COUNSELOR}</button>),
    ).toEqual(true);
    expect(wrapper.containsMatchingElement(<button type="button">Select {roleTypes.QAE_CSR}</button>)).toEqual(true);
    expect(wrapper.containsMatchingElement(<button type="button">Select {roleTypes.PPM}</button>)).toEqual(false);
  });

  it('handles setActiveRole with the selected role', () => {
    const mockSetActiveRole = jest.fn();
    const wrapper = mount(
      <SelectApplication
        userRoles={[{ roleType: roleTypes.TOO }, { roleType: roleTypes.TIO }]}
        setActiveRole={mockSetActiveRole}
      />,
    );

    const selectRoleButton = wrapper.find('button').first();
    selectRoleButton.simulate('click');
    expect(mockSetActiveRole).toHaveBeenCalledWith(roleTypes.TOO);
  });
});

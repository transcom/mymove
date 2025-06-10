import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';

import SelectApplication from './SelectApplication';

import store from 'shared/store';
import { roleTypes } from 'constants/userRoles';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => jest.fn(),
}));

describe('SelectApplication component', () => {
  it('renders the active role if one exists', () => {
    const mockSetActiveRole = jest.fn();
    const wrapper = mount(
      <Provider store={store}>
        <SelectApplication activeRole="myRole" userInactiveRoles={[]} setActiveRole={mockSetActiveRole} />
      </Provider>,
    );
    expect(wrapper.containsMatchingElement(<h2>Current role: myRole</h2>)).toEqual(true);
  });

  it('renders the first user role if there is no active role', () => {
    const mockSetActiveRole = jest.fn();
    const wrapper = mount(
      <Provider store={store}>
        <SelectApplication
          userInactiveRoles={[{ roleType: 'myFirstRole' }, { roleType: 'myOtherRole' }]}
          setActiveRole={mockSetActiveRole}
        />
        ,
      </Provider>,
    );
    expect(wrapper.containsMatchingElement(<h2>Current role: myFirstRole</h2>)).toEqual(true);
  });

  it('renders buttons for each of the user’s roles, and does not render buttons for roles the user doesn’t have', () => {
    const mockSetActiveRole = jest.fn();
    const wrapper = mount(
      <Provider store={store}>
        <SelectApplication
          userInactiveRoles={[
            { roleType: roleTypes.TOO },
            { roleType: roleTypes.TIO },
            { roleType: roleTypes.SERVICES_COUNSELOR },
            { roleType: roleTypes.QAE },
            { roleType: roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE },
          ]}
          setActiveRole={mockSetActiveRole}
        />
        ,
      </Provider>,
    );

    expect(wrapper.containsMatchingElement(<button type="button">Select {roleTypes.TOO}</button>)).toEqual(true);
    expect(wrapper.containsMatchingElement(<button type="button">Select {roleTypes.TIO}</button>)).toEqual(true);
    expect(
      wrapper.containsMatchingElement(<button type="button">Select {roleTypes.SERVICES_COUNSELOR}</button>),
    ).toEqual(true);
    expect(wrapper.containsMatchingElement(<button type="button">Select {roleTypes.QAE}</button>)).toEqual(true);
    expect(
      wrapper.containsMatchingElement(
        <button type="button">Select {roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE}</button>,
      ),
    ).toEqual(true);
  });

  it('handles setActiveRole with the selected role', () => {
    const mockSetActiveRole = jest.fn();
    const wrapper = mount(
      <Provider store={store}>
        <SelectApplication
          userInactiveRoles={[{ roleType: roleTypes.TOO }, { roleType: roleTypes.TIO }]}
          setActiveRole={mockSetActiveRole}
        />
        ,
      </Provider>,
    );

    const selectRoleButton = wrapper.find('button').first();
    selectRoleButton.simulate('click');
    expect(mockSetActiveRole).toHaveBeenCalledWith(roleTypes.TOO);
  });
});

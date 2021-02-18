import React from 'react';
import { mount } from 'enzyme';
import configureStore from 'redux-mock-store';

import LoginButton from './LoginButton';

describe('LoginButton tests', () => {
  const mockStore = configureStore();
  const initialState = {
    auth: {
      isLoggedIn: false,
    },
  };

  it('shows the EULA when the signin button is clicked and hides the EULA when cancel is clicked', () => {
    const store = mockStore(initialState);
    const wrapper = mount(<LoginButton store={store} />);
    expect(wrapper.find('[data-testid="modal"]').length).toEqual(0);
    wrapper.find('button[testid="signin"]').simulate('click');
    expect(wrapper.find('[data-testid="modal"]').length).toEqual(1);
    const CancelButton = wrapper.find('button[aria-label="Cancel"]');
    CancelButton.simulate('click');
    expect(wrapper.find('[data-testid="modal"]').length).toEqual(0);
  });

  it('renders the signin button when the user is not logged in', () => {
    const store = mockStore(initialState);
    const wrapper = mount(<LoginButton store={store} />);
    expect(wrapper.find('button[testid="signin"]').length).toEqual(1);
    expect(wrapper.find('a[testid="devlocal-signin"]').length).toEqual(0);
  });

  it('renders the devlocal signin button when running in development', () => {
    const store = mockStore({ ...initialState, isDevelopment: true });
    const wrapper = mount(<LoginButton store={store} />);
    expect(wrapper.find('button[testid="signin"]').length).toEqual(1);
    expect(wrapper.find('a[testid="devlocal-signin"]').length).toEqual(1);
  });
});

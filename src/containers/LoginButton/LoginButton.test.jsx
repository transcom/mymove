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

  it('renders the signin button when the user is not logged in', () => {
    const store = mockStore(initialState);
    const wrapper = mount(<LoginButton store={store} />);
    expect(wrapper.find('a[data-hook="signin"]').length).toEqual(1);
    expect(wrapper.find('a[data-hook="devlocal-signin"]').length).toEqual(0);
  });

  it('renders the devlocal signin button when running in development', () => {
    const store = mockStore({ ...initialState, isDevelopment: true });
    const wrapper = mount(<LoginButton store={store} />);
    expect(wrapper.find('a[data-hook="signin"]').length).toEqual(1);
    expect(wrapper.find('a[data-hook="devlocal-signin"]').length).toEqual(1);
  });
});

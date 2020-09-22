/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';

import Home from '.';

import { store } from 'shared/store';

const defaultProps = {
  serviceMember: {
    current_station: {},
  },
  showLoggedInUser: jest.fn(),
  createServiceMember: jest.fn(),
  loadPpms: jest.fn(),
  mtoShipment: {},
  isLoggedIn: true,
  loggedInUserIsLoading: false,
  loggedInUserSuccess: true,
  loggedInUserError: false,
  isProfileComplete: true,
  moveSubmitSuccess: false,
  ppm: {},
  loadMTOShipments: jest.fn(),
  orders: {},
  history: {},
  location: {},
  move: {},
};

function mountHome(props = defaultProps) {
  return mount(
    <Provider store={store}>
      <Home {...props} />
    </Provider>,
  );
}
describe('Home component', () => {
  it('renders Home with the right amount of components', () => {
    const wrapper = mountHome();
    expect(wrapper.find('Step').length).toBe(4);
    expect(wrapper.find('Helper').length).toBe(1);
    expect(wrapper.find('Contact').length).toBe(1);
  });
});

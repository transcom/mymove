/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import moment from 'moment';

import Home from '.';

import { store } from 'shared/store';

const defaultProps = {
  serviceMember: {
    current_station: {},
  },
  showLoggedInUser: jest.fn(),
  createServiceMember: jest.fn(),
  shipments: [],
  mtoShipment: {},
  isLoggedIn: true,
  loggedInUserIsLoading: false,
  loggedInUserSuccess: true,
  loggedInUserError: false,
  isProfileComplete: true,
  moveSubmitSuccess: false,
  currentPpm: {},
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
  describe('contents of Step 3', () => {
    it('contains ppm and hhg cards if those shipments exist', () => {
      let props = {
        currentPpm: { id: '12345', createdAt: moment() },
        shipments: [
          { id: '4321', createdAt: moment().add(1, 'days'), shipmentType: 'HHG' },
          { id: '4322', createdAt: moment().subtract(1, 'days'), shipmentType: 'HHG' },
        ],
      };

      props = { ...defaultProps, ...props };
      const wrapper = mountHome(props);
      expect(wrapper.find('ShipmentListItem').length).toBe(3);
    });
  });
});

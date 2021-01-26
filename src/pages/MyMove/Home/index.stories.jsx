/*  react/jsx-props-no-spreading */
import React from 'react';
import { Provider } from 'react-redux';

import Home from '.';

import { store } from 'shared/store';

export default {
  title: 'Customer Components / Home',
};

const defaultProps = {
  serviceMember: {
    first_name: 'John',
    last_name: 'Lee',
    current_station: {
      name: 'Fort Knox',
    },
  },
  showLoggedInUser() {},
  loadMTOShipments() {},
};
export const Basic = () => (
  <Provider store={store}>
    <div className="grid-container usa-prose">
      <Home {...defaultProps} />
    </div>
  </Provider>
);

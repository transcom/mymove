/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { Provider } from 'react-redux';

import Home from '.';

import { store } from 'shared/store';

export default {
  title: 'Customer Components | Home',
};

export const Basic = () => (
  <Provider store={store}>
    <div className="grid-container usa-prose">
      <Home showLoggedInUser={() => {}} />
    </div>
  </Provider>
);

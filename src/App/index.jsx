import React from 'react';
import { Provider } from 'react-redux';

import store from 'shared/store';
import './index.css';

import Loadable from 'react-loadable';

import { AppContext, officeContext, myMoveContext } from 'shared/AppContext';

const Office = Loadable({
  loader: () => import('scenes/Office'),
  loading: () => <div>Loading...</div>,
});

const MyMove = Loadable({
  loader: () => import('scenes/MyMove'),
  loading: () => <div>Loading...</div>,
});

const hostname = window && window.location && window.location.hostname;
const isOfficeSite = hostname.startsWith('office');
const App = () => {
  if (isOfficeSite)
    return (
      <Provider store={store}>
        <AppContext.Provider value={officeContext}>
          <Office />
        </AppContext.Provider>
      </Provider>
    );
  return (
    <Provider store={store}>
      <AppContext.Provider value={myMoveContext}>
        <MyMove />
      </AppContext.Provider>
    </Provider>
  );
};

export default App;

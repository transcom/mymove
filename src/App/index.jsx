import React from 'react';
import { Provider } from 'react-redux';
import { isOfficeSite, isTspSite } from 'shared/constants.js';
import { store } from 'shared/store';
import './index.css';

import Loadable from 'react-loadable';

import { AppContext, defaultTspContext, defaultOfficeContext, defaultMyMoveContext } from 'shared/AppContext';
import { detectFlags } from 'shared/featureFlags.js';

const Tsp = Loadable({
  loader: () => import('scenes/TransportationServiceProvider'),
  loading: () => <div>Loading...</div>,
});

const Office = Loadable({
  loader: () => import('scenes/Office'),
  loading: () => <div>Loading...</div>,
});

const MyMove = Loadable({
  loader: () => import('scenes/MyMove'),
  loading: () => <div>Loading...</div>,
});

const flags = detectFlags(process.env['NODE_ENV'], window.location.host, window.location.search);

const tspContext = Object.assign({}, defaultTspContext, { flags });
const officeContext = Object.assign({}, defaultOfficeContext, { flags });
const myMoveContext = Object.assign({}, defaultMyMoveContext, { flags });

const App = () => {
  if (isOfficeSite)
    return (
      <Provider store={store}>
        <AppContext.Provider value={officeContext}>
          <Office />
        </AppContext.Provider>
      </Provider>
    );
  else if (isTspSite)
    return (
      <Provider store={store}>
        <AppContext.Provider value={tspContext}>
          <Tsp />
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

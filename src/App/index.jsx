import React from 'react';
import { Provider } from 'react-redux';

import Loadable from 'react-loadable';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { isOfficeSite, isTspSite, isAdminSite, isSystemAdminSite } from 'shared/constants.js';
import { store } from 'shared/store';
import { AppContext, defaultOfficeContext, defaultMyMoveContext, defaultAdminContext } from 'shared/AppContext';
import { detectFlags } from 'shared/featureFlags.js';

import './index.css';

const Office = Loadable({
  loader: () => import('scenes/Office'),
  loading: () => <LoadingPlaceholder />,
});

const MyMove = Loadable({
  loader: () => import('scenes/MyMove'),
  loading: () => <LoadingPlaceholder />,
});

// Will uncomment for program admin
// const Admin = Loadable({
//   loader: () => import('scenes/Admin'),
//   loading: () => <LoadingPlaceholder />,
// });
//
const SystemAdmin = Loadable({
  loader: () => import('scenes/SystemAdmin'),
  loading: () => <LoadingPlaceholder />,
});

const flags = detectFlags(process.env['NODE_ENV'], window.location.host, window.location.search);

const officeContext = Object.assign({}, defaultOfficeContext, { flags });
const myMoveContext = Object.assign({}, defaultMyMoveContext, { flags });
const adminContext = Object.assign({}, defaultAdminContext, { flags });

const App = () => {
  if (isOfficeSite)
    return (
      <Provider store={store}>
        <AppContext.Provider value={officeContext}>
          <Office />
        </AppContext.Provider>
      </Provider>
    );
  else if (isTspSite) return <h1 style={{ textAlign: 'center' }}>TSP App</h1>;
  else if (isSystemAdminSite)
    return (
      <AppContext.Provider value={adminContext}>
        <SystemAdmin />
      </AppContext.Provider>
    );
  else if (isAdminSite) return <SystemAdmin />;
  return (
    <Provider store={store}>
      <AppContext.Provider value={myMoveContext}>
        <MyMove />
      </AppContext.Provider>
    </Provider>
  );
};

export default App;

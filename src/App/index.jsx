import React from 'react';
import { Provider } from 'react-redux';
import Loadable from 'react-loadable';
import { ConnectedRouter } from 'connected-react-router';
import { PersistGate } from 'redux-persist/integration/react';
import { ReactQueryConfigProvider } from 'react-query';
import { ReactQueryDevtools } from 'react-query-devtools';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { isOfficeSite, isAdminSite, isSystemAdminSite } from 'shared/constants';
import { store, persistor, history } from 'shared/store';
import { AppContext, defaultOfficeContext, defaultMyMoveContext, defaultAdminContext } from 'shared/AppContext';
import { detectFlags } from 'shared/featureFlags';

import '../icons';
import './index.css';

const Office = Loadable({
  loader: () => import('pages/Office'),
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

const flags = detectFlags(process.env.NODE_ENV, window.location.host, window.location.search);

const officeContext = { ...defaultOfficeContext, flags };
const myMoveContext = { ...defaultMyMoveContext, flags };
const adminContext = { ...defaultAdminContext, flags };

const officeQueryConfig = {
  queries: {
    retry: false, // default to no retries for now
    refetchOnWindowFocus: true,
    // onError: noop, // TODO - log errors?
  },
  mutations: {
    // onError: noop, // TODO - log errors?
  },
};

const App = () => {
  if (isOfficeSite)
    return (
      <ReactQueryConfigProvider config={officeQueryConfig}>
        <Provider store={store}>
          <PersistGate loading={<LoadingPlaceholder />} persistor={persistor}>
            <AppContext.Provider value={officeContext}>
              <ConnectedRouter history={history}>
                <Office />
                <ReactQueryDevtools initialIsOpen={false} />
              </ConnectedRouter>
            </AppContext.Provider>
          </PersistGate>
        </Provider>
      </ReactQueryConfigProvider>
    );

  if (isSystemAdminSite)
    return (
      <AppContext.Provider value={adminContext}>
        <SystemAdmin />
      </AppContext.Provider>
    );

  if (isAdminSite) return <SystemAdmin />;

  return (
    <Provider store={store}>
      <AppContext.Provider value={myMoveContext}>
        <MyMove />
      </AppContext.Provider>
    </Provider>
  );
};

export default App;

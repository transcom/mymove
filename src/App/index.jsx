import React, { lazy, Suspense } from 'react';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';
import { PersistGate } from 'redux-persist/integration/react';
import { ReactQueryConfigProvider } from 'react-query';
import { ReactQueryDevtools } from 'react-query-devtools';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { isOfficeSite, isAdminSite, isSystemAdminSite } from 'shared/constants';
import { store, persistor, history } from 'shared/store';
import { AppContext, defaultOfficeContext, defaultMyMoveContext, defaultAdminContext } from 'shared/AppContext';
import { detectFlags } from 'utils/featureFlags';
import '../icons';
import 'shared/shared.css';
import './index.css';
import MarkerIO from 'components/ThirdParty/MarkerIO';

const Office = lazy(() => import('pages/Office'));
const MyMove = lazy(() => import('scenes/MyMove'));

// Will uncomment for program admin
// const Admin = Loadable({
//   loader: () => import('scenes/Admin'),
//   loading: () => <LoadingPlaceholder />,
// });
//

const SystemAdmin = lazy(() => import('scenes/SystemAdmin'));

const flags = detectFlags(process.env.NODE_ENV, window.location.host, window.location.search);

const officeContext = { ...defaultOfficeContext, flags };
const myMoveContext = { ...defaultMyMoveContext, flags };
const adminContext = { ...defaultAdminContext, flags };

const officeQueryConfig = {
  queries: {
    retry: false, // default to no retries for now
    refetchOnWindowFocus: false,
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
                <Suspense fallback={<LoadingPlaceholder />}>
                  <Office />
                  {flags.markerIO && <MarkerIO />}
                </Suspense>
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
        <Suspense fallback={<LoadingPlaceholder />}>
          <SystemAdmin />
        </Suspense>
      </AppContext.Provider>
    );

  if (isAdminSite) return <SystemAdmin />;

  return (
    <Provider store={store}>
      <AppContext.Provider value={myMoveContext}>
        <ConnectedRouter history={history}>
          <Suspense fallback={<LoadingPlaceholder />}>
            <MyMove />
            {flags.markerIO && <MarkerIO />}
          </Suspense>
        </ConnectedRouter>
      </AppContext.Provider>
    </Provider>
  );
};

export default App;

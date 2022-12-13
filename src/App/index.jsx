import React, { lazy, Suspense } from 'react';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';
import { PersistGate } from 'redux-persist/integration/react';
import { QueryClientProvider, QueryClient } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { isOfficeSite, isAdminSite } from 'shared/constants';
import { store, persistor, history } from 'shared/store';
import { AppContext, defaultOfficeContext, defaultMyMoveContext, defaultAdminContext } from 'shared/AppContext';
import { detectFlags } from 'utils/featureFlags';
import '../icons';
import 'shared/shared.css';
import './index.css';
import MarkerIO from 'components/ThirdParty/MarkerIO';
import ScrollToTop from 'components/ScrollToTop';

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

const officeQueryConfig = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false, // default to no retries for now
      // do not re-query on window refocus
      refetchOnWindowFocus: false,
      // onError: noop, // TODO - log errors?
      networkMode: 'offlineFirst', // restoring previous-behavior. Without this, it will be paused without a network
    },
    mutations: {
      // onError: noop, // TODO - log errors?
      networkMode: 'offlineFirst', // restoring previous-behavior. Without this, it will be paused without a network
    },
  },
});

const App = () => {
  if (isOfficeSite)
    return (
      <QueryClientProvider config={officeQueryConfig}>
        <Provider store={store}>
          <PersistGate loading={<LoadingPlaceholder />} persistor={persistor}>
            <AppContext.Provider value={officeContext}>
              <ConnectedRouter history={history}>
                <Suspense fallback={<LoadingPlaceholder />}>
                  <ScrollToTop />
                  <Office />
                  {flags.markerIO && <MarkerIO />}
                </Suspense>
                <ReactQueryDevtools initialIsOpen={false} />
              </ConnectedRouter>
            </AppContext.Provider>
          </PersistGate>
        </Provider>
      </QueryClientProvider>
    );

  if (isAdminSite)
    return (
      <Provider store={store}>
        <PersistGate loading={<LoadingPlaceholder />} persistor={persistor}>
          <AppContext.Provider value={adminContext}>
            <ConnectedRouter history={history}>
              <Suspense fallback={<LoadingPlaceholder />}>
                <SystemAdmin />
              </Suspense>
            </ConnectedRouter>
          </AppContext.Provider>
        </PersistGate>
      </Provider>
    );

  return (
    <Provider store={store}>
      <AppContext.Provider value={myMoveContext}>
        <ConnectedRouter history={history}>
          <Suspense fallback={<LoadingPlaceholder />}>
            <ScrollToTop />
            <MyMove />
            {flags.markerIO && <MarkerIO />}
          </Suspense>
        </ConnectedRouter>
      </AppContext.Provider>
    </Provider>
  );
};

export default App;

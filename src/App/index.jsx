import React, { lazy, Suspense } from 'react';
import { Provider } from 'react-redux';
import { BrowserRouter } from 'react-router-dom';
import { PersistGate } from 'redux-persist/integration/react';
import { QueryClientProvider, QueryClient } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { isOfficeSite, isAdminSite } from 'shared/constants';
import { store, persistor } from 'shared/store';
import { AppContext, defaultOfficeContext, defaultMyMoveContext, defaultAdminContext } from 'shared/AppContext';
import { detectFlags } from 'utils/featureFlags';
import '../icons';
import 'shared/shared.css';
import './index.css';
import MarkerIO from 'components/ThirdParty/MarkerIO';
import MilMoveErrorBoundary from 'components/MilMoveErrorBoundary';
import ScrollToTop from 'components/ScrollToTop';
import PageTitle from 'components/PageTitle';
import MaintenancePage from 'pages/Maintenance/MaintenancePage';

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
  // if (flags.InMaintenance === undefined) {
  //   return <MaintenancePage />;
  // }

  // We need an error boundary around each of the main apps (Office,
  // SystemAdmin, MyMove) because they are lazy loaded and it's
  // possible we could get a ChunkLoadError when trying to load them.
  // Each of the main apps has its own componentDidCatch which would
  // mean the MilMoveErrorBoundary is probably unlikely to be reached
  if (isOfficeSite) {
    return (
      <MilMoveErrorBoundary fallback={<SomethingWentWrong />}>
        <QueryClientProvider client={officeQueryConfig}>
          <Provider store={store}>
            <PersistGate loading={<LoadingPlaceholder />} persistor={persistor}>
              <AppContext.Provider value={officeContext}>
                <BrowserRouter>
                  <Suspense fallback={<LoadingPlaceholder />}>
                    <ScrollToTop />
                    <PageTitle />
                    <Office />
                    {flags.markerIO && <MarkerIO />}
                  </Suspense>
                  <ReactQueryDevtools initialIsOpen={false} />
                </BrowserRouter>
              </AppContext.Provider>
            </PersistGate>
          </Provider>
        </QueryClientProvider>
      </MilMoveErrorBoundary>
    );
  }
  if (isAdminSite) {
    return (
      <MilMoveErrorBoundary fallback={<SomethingWentWrong />}>
        <Provider store={store}>
          <AppContext.Provider value={adminContext}>
            <BrowserRouter>
              <Suspense fallback={<LoadingPlaceholder />}>
                <PageTitle />
                <SystemAdmin />
              </Suspense>
            </BrowserRouter>
          </AppContext.Provider>
        </Provider>
      </MilMoveErrorBoundary>
    );
  }

  return (
    <MilMoveErrorBoundary fallback={<SomethingWentWrong />}>
      <Provider store={store}>
        <AppContext.Provider value={myMoveContext}>
          <BrowserRouter>
            <Suspense fallback={<LoadingPlaceholder />}>
              <ScrollToTop />
              <PageTitle />
              <MyMove />
              {flags.markerIO && <MarkerIO />}
            </Suspense>
          </BrowserRouter>
        </AppContext.Provider>
      </Provider>
    </MilMoveErrorBoundary>
  );
};

export default App;

import React, { Suspense } from 'react';
import { node, shape, arrayOf, func, string } from 'prop-types';
import { Provider } from 'react-redux';
import { createMemoryHistory } from 'history';
import { ConnectedRouter } from 'connected-react-router';
import { Router } from 'react-router-dom';
/* eslint-disable-next-line import/no-extraneous-dependencies */
import { render } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import { configureStore } from 'shared/store';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import PermissionProvider from 'components/Restricted/PermissionProvider';

export const createMockHistory = (initialEntries) => {
  return createMemoryHistory({ initialEntries });
};

export const ReactQueryWrapper = ({ children, client }) => {
  const queryClient = client || new QueryClient();
  return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
};

export const renderWithRouter = (ui, { route = '/', history = createMockHistory([route]) } = {}) => {
  return {
    ...render(
      <ReactQueryWrapper>
        <Router history={history}>{ui}</Router>
      </ReactQueryWrapper>,
    ),
    history,
  };
};

export const MockProviders = ({
  children,
  initialState,
  initialEntries,
  permissions,
  history,
  currentUserId,
  client,
}) => {
  const mockHistory = history || createMockHistory(initialEntries);
  const mockStore = configureStore(mockHistory, initialState);

  return (
    <ReactQueryWrapper client={client}>
      <PermissionProvider permissions={permissions} currentUserId={currentUserId}>
        <Provider store={mockStore.store}>
          <ConnectedRouter history={mockHistory}>
            <Suspense fallback={<LoadingPlaceholder />}>{children}</Suspense>
          </ConnectedRouter>
        </Provider>
      </PermissionProvider>
    </ReactQueryWrapper>
  );
};

MockProviders.propTypes = {
  children: node.isRequired,
  initialState: shape({}),
  initialEntries: arrayOf(string),
  history: shape({
    push: func.isRequired,
    goBack: func.isRequired,
  }),
  permissions: arrayOf(string),
  currentUserId: string,
};

const DEFAULT_INITIAL_ENTRIES = ['/'];

MockProviders.defaultProps = {
  initialState: {},
  initialEntries: DEFAULT_INITIAL_ENTRIES,
  history: null,
  permissions: [],
  currentUserId: null,
};

export const setUpProvidersWithHistory = (initialEntries = DEFAULT_INITIAL_ENTRIES) => {
  const memoryHistory = createMockHistory(initialEntries);

  const mockProviderWithHistory = (props) => <MockProviders history={memoryHistory} {...props} />;

  return {
    memoryHistory,
    mockProviderWithHistory,
  };
};

import React, { Suspense } from 'react';
import { node, shape, arrayOf, func, string } from 'prop-types';
import { Provider } from 'react-redux';
import { createMemoryHistory } from 'history';
import { ConnectedRouter } from 'connected-react-router';
import { Router } from 'react-router-dom';
/* eslint-disable-next-line import/no-extraneous-dependencies */
import { render } from '@testing-library/react';

import { configureStore } from 'shared/store';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import PermissionProvider from 'components/Restricted/PermissionProvider';

export const createMockHistory = (initialEntries) => {
  return createMemoryHistory({ initialEntries });
};

export const renderWithRouter = (ui, { route = '/', history = createMockHistory([route]) } = {}) => {
  return {
    ...render(<Router history={history}>{ui}</Router>),
    history,
  };
};

export const MockProviders = ({ children, initialState, initialEntries, permissions, history, currentUserId }) => {
  const mockHistory = history || createMockHistory(initialEntries);
  const mockStore = configureStore(mockHistory, initialState);

  return (
    <PermissionProvider permissions={permissions} currentUserId={currentUserId}>
      <Provider store={mockStore.store}>
        <ConnectedRouter history={mockHistory}>
          <Suspense fallback={<LoadingPlaceholder />}>{children}</Suspense>
        </ConnectedRouter>
      </Provider>
    </PermissionProvider>
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

MockProviders.defaultProps = {
  initialState: {},
  initialEntries: ['/'],
  history: null,
  permissions: [],
  currentUserId: null,
};

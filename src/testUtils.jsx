import React, { Suspense } from 'react';
import { node, shape, arrayOf, func, string } from 'prop-types';
import { Provider } from 'react-redux';
import { createMemoryHistory } from 'history';
import { ConnectedRouter } from 'connected-react-router';
import { Router } from 'react-router-dom';
import { render } from '@testing-library/react'; // eslint-disable-line import/no-extraneous-dependencies

import { configureStore } from 'shared/store';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

export const createMockHistory = (initialEntries) => {
  return createMemoryHistory({ initialEntries });
};

export const renderWithRouter = (ui, { route = '/', history = createMockHistory([route]) } = {}) => {
  return {
    ...render(<Router history={history}>{ui}</Router>),
    history,
  };
};

export const MockProviders = ({ children, initialState, initialEntries, history }) => {
  const mockHistory = history || createMockHistory(initialEntries);
  const mockStore = configureStore(mockHistory, initialState);

  return (
    <Provider store={mockStore.store}>
      <ConnectedRouter history={mockHistory}>
        <Suspense fallback={<LoadingPlaceholder />}>{children}</Suspense>
      </ConnectedRouter>
    </Provider>
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
};

MockProviders.defaultProps = {
  initialState: {},
  initialEntries: ['/'],
  history: null,
};

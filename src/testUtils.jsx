import React from 'react';
import { node, shape, arrayOf, func, string } from 'prop-types';
import { Provider } from 'react-redux';
import { createMemoryHistory } from 'history';
import { ConnectedRouter } from 'connected-react-router';

import { configureStore } from 'shared/store';

export const createMockHistory = (initialEntries) => {
  return createMemoryHistory({ initialEntries });
};

export const MockProviders = ({ children, initialState, initialEntries, history }) => {
  const mockHistory = history || createMockHistory(initialEntries);
  const mockStore = configureStore(mockHistory, initialState);

  return (
    <Provider store={mockStore.store}>
      <ConnectedRouter history={mockHistory}>{children}</ConnectedRouter>
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

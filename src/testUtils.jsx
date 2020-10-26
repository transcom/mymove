/*  import/prefer-default-export */
import React from 'react';
import PropTypes from 'prop-types';
import { Provider } from 'react-redux';
import { createMemoryHistory } from 'history';
import { ConnectedRouter } from 'connected-react-router';

import { configureStore } from 'shared/store';

export const MockProviders = ({ children, initialState = {}, initialEntries = [] }) => {
  const mockHistory = createMemoryHistory({ initialEntries });
  const mockStore = configureStore(mockHistory, initialState);

  return (
    <Provider store={mockStore.store}>
      <ConnectedRouter history={mockHistory}>{children}</ConnectedRouter>
    </Provider>
  );
};

MockProviders.propTypes = {
  children: PropTypes.node.isRequired,
  //  react/forbid-prop-types
  initialState: PropTypes.object,
  initialEntries: PropTypes.arrayOf(PropTypes.string),
};

MockProviders.defaultProps = {
  initialState: {},
  initialEntries: [],
};

import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import MockRouter from 'react-mock-router';
import { schema } from 'normalizr';

import MoveInfo from './MoveInfo';
import store from 'shared/store';

const dummyFunc = () => {};
const loadDependenciesHasError = null;
const loadDependenciesHasSuccess = false;
const location = {
  pathname: '',
};
const match = {
  params: { moveID: '123456' },
  url: 'www.nino.com',
  path: '/moveIt/moveIt',
};
const mockStore = {
  ...store,
  entities: {
    moves: {},
  },
  schema,
};

const push = jest.fn();

describe('Loads MoveInfo', () => {
  it.skip('renders without crashing', () => {
    const div = document.createElement('div');
    ReactDOM.render(
      <Provider store={mockStore}>
        <MockRouter push={push}>
          <MoveInfo
            loadDependenciesHasError={loadDependenciesHasError}
            loadDependenciesHasSuccess={loadDependenciesHasSuccess}
            location={location}
            match={match}
            loadMoveDependencies={dummyFunc}
          />
        </MockRouter>
      </Provider>,
      div,
    );
  });
});

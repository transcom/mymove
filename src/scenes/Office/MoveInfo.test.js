import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import MockRouter from 'react-mock-router';

import MoveInfo from './MoveInfo';
import store from 'shared/store';

const dummyFunc = () => {};
const moveIsLoading = false;
const ordersAreLoading = false;
const serviceMemberIsLoading = false;
const backupContactsAreLoading = false;
const moveHasLoadError = false;
const moveHasLoadSuccess = null;
const ordersHaveLoadError = false;
const ordersHaveLoadSuccess = null;
const serviceMemberHasLoadError = false;
const serviceMemberHasLoadSuccess = null;
const backupContactsHaveLoadError = false;
const backupContactsHaveLoadSuccess = null;
const match = {
  params: { moveID: '123456' },
  url: 'www.nino.com',
  path: '/moveIt/moveIt',
};

const push = jest.fn();

describe('Loads MoveInfo', () => {
  it('renders without crashing', () => {
    const div = document.createElement('div');
    ReactDOM.render(
      <Provider store={store}>
        <MockRouter push={push}>
          <MoveInfo
            moveIsLoading={moveIsLoading}
            moveHasLoadError={moveHasLoadError}
            moveHasLoadSuccess={moveHasLoadSuccess}
            serviceMemberIsLoading={serviceMemberIsLoading}
            serviceMemberHasLoadError={serviceMemberHasLoadError}
            serviceMemberHasLoadSuccess={serviceMemberHasLoadSuccess}
            ordersAreLoading={ordersAreLoading}
            ordersHaveLoadError={ordersHaveLoadError}
            ordersHaveLoadSuccess={ordersHaveLoadSuccess}
            backupContactsAreLoading={backupContactsAreLoading}
            backupContactsHaveLoadError={backupContactsHaveLoadError}
            backupContactsHaveLoadSuccess={backupContactsHaveLoadSuccess}
            match={match}
            loadMoveDependencies={dummyFunc}
          />
        </MockRouter>
      </Provider>,
      div,
    );
  });
});

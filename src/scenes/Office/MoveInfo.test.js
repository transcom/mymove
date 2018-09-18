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
const moveHasLoadError = null;
const moveHasLoadSuccess = false;
const ordersHaveLoadError = null;
const ordersHaveLoadSuccess = false;
const serviceMemberHasLoadError = null;
const serviceMemberHasLoadSuccess = false;
const backupContactsHaveLoadError = null;
const backupContactsHaveLoadSuccess = false;
const PPMsHaveLoadError = null;
const PPMsHaveLoadSuccess = false;
const loadDependenciesHasError = null;
const loadDependenciesHasSuccess = false;
const moveIsCanceling = false;
const moveHasCancelError = null;
const moveHasCancelSuccess = false;
const location = {
  pathname: '',
};
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
            PPMsHaveLoadError={PPMsHaveLoadError}
            PPMsHaveLoadSuccess={PPMsHaveLoadSuccess}
            loadDependenciesHasError={loadDependenciesHasError}
            loadDependenciesHasSuccess={loadDependenciesHasSuccess}
            location={location}
            match={match}
            loadMoveDependencies={dummyFunc}
            moveIsCanceling={moveIsCanceling}
            moveHasCancelError={moveHasCancelError}
            moveHasCancelSuccess={moveHasCancelSuccess}
          />
        </MockRouter>
      </Provider>,
      div,
    );
  });
});

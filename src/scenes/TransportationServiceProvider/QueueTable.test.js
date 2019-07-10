import React from 'react';
import { Provider } from 'react-redux';
import MockRouter from 'react-mock-router';
import configureMockStore from 'redux-mock-store';
import thunk from 'redux-thunk';

import QueueTable from './QueueTable';
import store from 'shared/store';
import { mount } from 'enzyme/build';
import { setIsLoggedInType } from 'shared/Data/users';

const push = jest.fn();

describe('Refreshing', () => {
  it('loads the data again', done => {
    const refreshSpy = jest.spyOn(QueueTable.WrappedComponent.prototype, 'refresh');
    const fetchDataSpy = jest.spyOn(QueueTable.WrappedComponent.prototype, 'fetchData');

    const wrapper = mountComponents(retrieveShipmentsStub());

    wrapper
      .find('[data-cy="refreshQueue"]')
      .at(0)
      .simulate('click');

    setTimeout(() => {
      expect(refreshSpy).toHaveBeenCalled();
      expect(fetchDataSpy).toHaveBeenCalled();

      done();
    });
  });
});

describe('on 401 unauthorized error', () => {
  const middlewares = [thunk];
  const mockStore = configureMockStore(middlewares);

  it('force user log out', done => {
    let error = new Error('Unauthorized');
    error.status = 401;

    const store = mockStore({});
    const wrapper = mountComponents(retrieveShipmentsStub(null, error), 'new', store);
    wrapper
      .find('[data-cy="refreshQueue"]')
      .at(0)
      .simulate('click');

    setTimeout(() => {
      const userLoggedOutAction = { type: setIsLoggedInType, isLoggedIn: false };
      expect(store.getActions()).toContainEqual(userLoggedOutAction);

      done();
    });
  });
});

function retrieveShipmentsStub(params, throwError) {
  // This is meant as a stub that will act in place of
  // `RetrieveShipmentsForTSP` from TransporationServiceProvider/api.js
  return async () => {
    if (throwError) {
      throw throwError;
    }

    return await new Promise(resolve => {
      resolve([
        {
          id: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
          status: '',
          service_member: {},
          traffic_distribution_list: {},
          ...params,
        },
      ]);
    });
  };
}

function mountComponents(getShipments, queueType = 'new', mockStore = store) {
  return mount(
    <Provider store={mockStore}>
      <MockRouter push={push}>
        <QueueTable queueType={queueType} retrieveShipments={getShipments} />
      </MockRouter>
    </Provider>,
  );
}

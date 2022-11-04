import React from 'react';
import { Provider } from 'react-redux';
import { HashRouter as Router } from 'react-router-dom';

import configureMockStore from 'redux-mock-store';
import thunk from 'redux-thunk';

import QueueTable from './QueueTable';
import ReactTable from 'react-table';
import store from 'shared/store';
import { mount } from 'enzyme/build';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const push = jest.fn();

describe('Shipments column', () => {
  let wrapper;

  it('renders "PPM" when it is a PPM move', (done) => {
    wrapper = mountComponents(
      retrieveMovesStub({
        ppm_status: 'PAYMENT_REQUESTED',
        hhg_status: undefined,
      }),
    );

    setTimeout(() => {
      const move = getMove(wrapper);
      expect(move.shipments).toEqual(SHIPMENT_OPTIONS.PPM);

      done();
    });
  });

  it('does not display when the queue type is anything other than "new"', (done) => {
    wrapper = mountComponents(
      retrieveMovesStub({
        ppm_status: undefined,
        hhg_status: undefined,
      }),
      'ppm',
    );

    setTimeout(() => {
      const move = getMove(wrapper);
      expect(move.shipments);

      done();
    });
  });
});

describe('Refreshing', () => {
  let wrapper;
  it('loads the data again', (done) => {
    const refreshSpy = jest.spyOn(QueueTable.WrappedComponent.prototype, 'refresh');
    const fetchDataSpy = jest.spyOn(QueueTable.WrappedComponent.prototype, 'fetchData');

    wrapper = mountComponents(retrieveMovesStub());

    wrapper.find('[data-testid="refreshQueue"]').at(0).simulate('click');

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

  it('force user log out', (done) => {
    const fetchDataSpy = jest.spyOn(QueueTable.WrappedComponent.prototype, 'fetchData');

    let error = new Error('Unauthorized');
    error.status = 401;

    const store = mockStore({});
    const wrapper = mountComponents(retrieveMovesStub(null, error), 'new', store);
    wrapper.find('[data-testid="refreshQueue"]').at(0).simulate('click');

    setTimeout(() => {
      expect(fetchDataSpy).toHaveBeenCalled();

      const logOut = { type: 'LOG_OUT' };
      expect(store.getActions()).toContainEqual(logOut);

      done();
    });
  });
});

function retrieveMovesStub(params, throwError) {
  // This is meant as a stub that will act in place of
  // `RetrieveMovesForOffice` from Office/api.js
  return async () => {
    return await new Promise((resolve) => {
      if (throwError) {
        throw throwError;
      }

      resolve([
        {
          id: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
          ...params,
        },
      ]);
    });
  };
}

function mountComponents(getMoves, queueType = 'new', mockStore = store) {
  return mount(
    <Provider store={mockStore}>
      <Router push={push}>
        <QueueTable queueType={queueType} retrieveMoves={getMoves} />
      </Router>
    </Provider>,
  );
}

function getMove(wrapper) {
  return wrapper.find(ReactTable).state().data[0];
}

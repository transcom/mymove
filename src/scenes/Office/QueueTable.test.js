import React from 'react';
import { Provider } from 'react-redux';
import MockRouter from 'react-mock-router';

import QueueTable from './QueueTable';
import ReactTable from 'react-table';
import store from 'shared/store';
import { mount } from 'enzyme/build';

const push = jest.fn();

describe('Shipments column', () => {
  let wrapper;

  it('renders "PPM" when it is a PPM move', done => {
    wrapper = mountComponents(
      retrieveMovesStub({
        ppm_status: 'PAYMENT_REQUESTED',
        hhg_status: undefined,
      }),
    );

    setTimeout(() => {
      const move = getMove(wrapper);
      expect(move.shipments).toEqual('PPM');

      done();
    });
  });

  it('renders "HHG" when it is a HHG move', done => {
    wrapper = mountComponents(
      retrieveMovesStub({
        ppm_status: undefined,
        hhg_status: 'APPROVED',
      }),
    );

    setTimeout(() => {
      const move = getMove(wrapper);
      expect(move.shipments).toEqual('HHG');

      done();
    });
  });

  it('renders "HHG, PPM" when it is a combo move', done => {
    wrapper = mountComponents(
      retrieveMovesStub({
        ppm_status: 'PAYMENT_REQUESTED',
        hhg_status: 'APPROVED',
      }),
    );

    setTimeout(() => {
      const move = getMove(wrapper);
      expect(move.shipments).toEqual('HHG, PPM');

      done();
    });
  });

  it('does not display when the queue type is anything other than "new"', done => {
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
  it('loads the data again', done => {
    const refreshSpy = jest.spyOn(QueueTable.WrappedComponent.prototype, 'refresh');
    const fetchDataSpy = jest.spyOn(QueueTable.WrappedComponent.prototype, 'fetchData');

    wrapper = mountComponents(retrieveMovesStub());

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

function retrieveMovesStub(params) {
  // This is meant as a stub that will act in place of
  // `RetrieveMovesForOffice` from Office/api.js
  return async () => {
    return await new Promise(resolve => {
      resolve([
        {
          id: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
          ...params,
        },
      ]);
    });
  };
}

function mountComponents(getMoves, queueType = 'new') {
  return mount(
    <Provider store={store}>
      <MockRouter push={push}>
        <QueueTable queueType={queueType} retrieveMoves={getMoves} />
      </MockRouter>
    </Provider>,
  );
}

function getMove(wrapper) {
  return wrapper.find(ReactTable).state().data[0];
}

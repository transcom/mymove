import React from 'react';
import { Provider } from 'react-redux';
import MockRouter from 'react-mock-router';

import QueueTable from './QueueTable';
import store from 'shared/store';
import { mount } from 'enzyme/build';

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

function retrieveShipmentsStub(params) {
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

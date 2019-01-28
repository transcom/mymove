import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';

import { StorageInTransitPanel } from './StorageInTransitPanel';

const mockStore = configureStore();
let store;

describe('StorageInTransit tests', () => {
  describe('When no items exist', () => {
    let wrapper;
    const sitRequests = [];

    store = mockStore({});
    wrapper = mount(
      <Provider store={store}>
        <StorageInTransitPanel sitRequests={sitRequests} shipmentId="" sitEntitlement={90} />
      </Provider>,
    );

    it('renders without crashing', () => {
      expect(wrapper.find('.storage-in-transit-panel').length).toEqual(1);
    });
  });
});

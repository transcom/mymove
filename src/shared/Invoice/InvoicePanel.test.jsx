import React from 'react';
import { shallow, mount } from 'enzyme';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';

import { InvoicePanel } from './InvoicePanel';

describe('InvoicePanel tests', () => {
  let wrapper;
  const shipmentLineItems = [''];
  const mockStore = configureStore();
  let store;
  beforeEach(() => {
    store = mockStore({});

    wrapper = mount(
      <Provider store={store}>
        <InvoicePanel shipmentLineItems={shipmentLineItems} />
      </Provider>,
    );
  });

  describe('When no items exist', () => {
    it('renders without crashing', () => {
      expect(wrapper.find('.empty-content').length).toEqual(1);
    });
  });
});

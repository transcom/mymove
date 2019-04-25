import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';
import AdminWrapper from './index';

let store;
const mockStore = configureStore();

describe('AdminIndex tests', () => {
  describe('AdminIndex home page', () => {
    let wrapper;
    store = mockStore({});
    wrapper = mount(
      <Provider store={store}>
        <AdminWrapper />
      </Provider>,
    );

    it('renders without crashing', () => {
      expect(wrapper.find('.admin-system-wrapper').length).toEqual(1);
    });
  });
});

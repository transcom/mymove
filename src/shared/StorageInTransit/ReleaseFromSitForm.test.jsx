import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';
import { ReleaseFromSitForm } from './ReleaseFromSitForm';

let store;
const mockStore = configureStore();
const submit = jest.fn();

const storageInTransitSchema = {
  properties: {
    released_on: {
      type: 'string',
      format: 'date',
      example: '2018-04-26',
      title: 'Released on',
    },
  },
};

describe('ReleaseFromSitForm tests', () => {
  describe('Pre-filled form', () => {
    let wrapper;
    store = mockStore({});
    wrapper = mount(
      <Provider store={store}>
        <ReleaseFromSitForm onSubmit={submit} storageInTransitSchema={storageInTransitSchema} />
      </Provider>,
    );

    it('renders without crashing', () => {
      expect(wrapper.find('.release-from-sit-form').length).toEqual(1);
    });
  });
});

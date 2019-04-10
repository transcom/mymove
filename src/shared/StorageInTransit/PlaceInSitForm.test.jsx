import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';
import { PlaceInSitForm } from './PlaceInSitForm';

let store;
const mockStore = configureStore();
const submit = jest.fn();

const storageInTransitSchema = {
  properties: {
    actual_start_date: {
      type: 'string',
      format: 'date',
      example: '2018-04-26',
      title: 'Actual start date',
    },
  },
};

describe('PlaceInSitForm tests', () => {
  describe('Pre-filled form', () => {
    let wrapper;
    store = mockStore({});
    wrapper = mount(
      <Provider store={store}>
        <PlaceInSitForm onSubmit={submit} storageInTransitSchema={storageInTransitSchema} />
      </Provider>,
    );

    it('renders without crashing', () => {
      expect(wrapper.find('.place-in-sit-form').length).toEqual(1);
    });
  });
});

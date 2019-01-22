import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';

import { PreApprovalPanel } from './PreApprovalPanel';

describe('PreApprovalPanel tests', () => {
  let wrapper;
  const shipmentLineItems = [
    {
      id: 'sldkjf',
      tariff400ng_item: { code: '105D', item: 'Reg Shipping' },
      location: 'D',
      base_quantity: ' 16.7',
      notes: '',
      created_at: '2018-09-24T14:05:38.847Z',
      status: 'SUBMITTED',
    },
    {
      id: 'sldsdff',
      tariff400ng_item: { code: '105D', item: 'Reg Shipping' },
      location: 'D',
      base_quantity: ' 16.7',
      notes: 'Mounted deer head measures 23" x 34" x 27"; crate will be 16.7 cu ft',
      created_at: '2018-09-24T14:05:38.847Z',
      status: 'APPROVED',
    },
  ];
  const tariff400ng_items = [
    {
      id: 'sdlfkj',
      code: 'F9D',
      item: 'Long Haul',
    },
    {
      id: 'badfka',
      code: '19D',
      item: 'Crate',
    },
  ];
  const mockStore = configureStore();
  let store;
  beforeEach(() => {
    store = mockStore({});
    //mount appears to be necessary to get inner components to load (i.e. tests fail with shallow)
    wrapper = mount(
      <Provider store={store}>
        <PreApprovalPanel shipmentLineItems={shipmentLineItems} tariff400ng_items={tariff400ng_items} />
      </Provider>,
    );
  });

  describe('When on approval is passed in and status is submitted', () => {
    it('renders without crashing', () => {
      const icons = wrapper.find('.icon');
      expect(wrapper.find('.pre-approval-panel').length).toEqual(1);
      expect(icons.length).toBe(8);
    });
  });
  describe('When creator and request sub-components are activated', () => {
    it('a request activation hides the creator', () => {
      expect(wrapper.find('Creator').length).toBe(1);
      wrapper
        .find('PreApprovalPanel')
        .instance()
        .onRequestActivation(true);
      wrapper.update();
      expect(wrapper.find('Creator').length).toBe(0);
      wrapper
        .find('PreApprovalPanel')
        .instance()
        .onRequestActivation(false);
      wrapper.update();
      expect(wrapper.find('Creator').length).toBe(1);
    });
    it('a creator activation deactivates the table', () => {
      expect(
        wrapper
          .find('PreApprovalTable')
          .first()
          .prop('isActionable'),
      ).toBe(true);
      wrapper
        .find('PreApprovalPanel')
        .instance()
        .onFormActivation(true);
      wrapper.update();
      expect(
        wrapper
          .find('PreApprovalTable')
          .first()
          .prop('isActionable'),
      ).toBe(false);
      wrapper
        .find('PreApprovalPanel')
        .instance()
        .onFormActivation(false);
      wrapper.update();
      expect(
        wrapper
          .find('PreApprovalTable')
          .first()
          .prop('isActionable'),
      ).toBe(true);
    });
  });
});

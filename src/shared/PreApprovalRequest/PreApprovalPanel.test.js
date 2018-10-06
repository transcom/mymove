import React from 'react';
import { shallow } from 'enzyme';
import { Provider } from 'react-redux';

import configureStore from 'redux-mock-store';
import { mount } from 'enzyme';

import PreApprovalPanel from './PreApprovalPanel.jsx';

describe('PreApprovalPanel tests', () => {
  let wrapper, icons;
  const onEdit = jest.fn();
  const shipment_accessorials = [
    {
      code: '105D',
      item: 'Unpack Reg Crate',
      location: 'D',
      base_quantity: ' 16.7',
      notes: '',
      created_at: '2018-09-24T14:05:38.847Z',
      status: 'SUBMITTED',
    },
    {
      code: '105E',
      item: 'Unpack Reg Crate',
      location: 'D',
      base_quantity: ' 16.7',
      notes:
        'Mounted deer head measures 23" x 34" x 27"; crate will be 16.7 cu ft',
      created_at: '2018-09-24T14:05:38.847Z',
      status: 'APPROVED',
    },
  ];
  const accessorials = [
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
        <PreApprovalPanel
          shipment_accessorials={shipment_accessorials}
          accessorials={accessorials}
        />
      </Provider>,
    );
  });

  describe('When on approval is passed in and status is submitted', () => {
    it('renders without crashing', () => {
      const childContainer = wrapper.find('.icon');
      expect(wrapper.find('.accessorial-panel').length).toEqual(1);
      expect(childContainer.length).toBe(6);
    });
  });
});

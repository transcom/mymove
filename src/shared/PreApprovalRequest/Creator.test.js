import React from 'react';
import { Provider } from 'react-redux';

import configureStore from 'redux-mock-store';
import { mount } from 'enzyme';

import { Creator } from './Creator';
import { no_op } from 'shared/utils';

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
const submit = jest.fn();
const clear = jest.fn();
const mockStore = configureStore();
let store;
let wrapper;
describe('given a Creator', () => {
  describe('when the form is disabled', () => {
    beforeEach(() => {
      store = mockStore({});
      //mount appears to be necessary to get inner components to load (i.e. tests fail with shallow)
      wrapper = mount(
        <Provider store={store}>
          <Creator
            accessorials={accessorials}
            submitForm={submit}
            formEnabled={false}
            hasSubmitSucceeded={false}
            savePreApprovalRequest={no_op}
            clearForm={clear}
          />
        </Provider>,
      );
    });

    it('renders without crashing', () => {
      // eslint-disable-next-line
      expect(wrapper.exists('a')).toBe(true);
    });
    it('can toggle to show the form', () => {
      wrapper.find('a').simulate('click');
      expect(wrapper.exists('div.accessorial-panel-modal')).toBe(true);
    });
    it('cancel closes the form', () => {
      wrapper.find('a').simulate('click');
      expect(wrapper.exists('div.accessorial-panel-modal')).toBe(true);
      wrapper.find('a').simulate('click');
      expect(wrapper.exists('div.accessorial-panel-modal')).toBe(false);
    });
    it('buttons are disabled', () => {
      wrapper.find('a').simulate('click');
      expect(
        wrapper.find('button.button-primary').prop('disabled'),
      ).toBeTruthy();
      expect(
        wrapper.find('button.button-secondary').prop('disabled'),
      ).toBeTruthy();
    });
  });
  describe('when the form is enabled', () => {
    beforeEach(() => {
      submit.mockClear();
      store = mockStore({});
      //mount appears to be necessary to get inner components to load (i.e. tests fail with shallow)
      wrapper = mount(
        <Provider store={store}>
          <Creator
            accessorials={accessorials}
            submitForm={submit}
            formEnabled={true}
            hasSubmitSucceeded={false}
            savePreApprovalRequest={no_op}
            clearForm={clear}
          />
        </Provider>,
      );
    });

    it('renders without crashing', () => {
      // eslint-disable-next-line
      expect(wrapper.exists('a')).toBe(true);
    });
    it('can toggle to show the form', () => {
      wrapper.find('a').simulate('click');
      expect(wrapper.exists('div.accessorial-panel-modal')).toBe(true);
    });
    it('cancel closes the form', () => {
      wrapper.find('a').simulate('click');
      expect(wrapper.exists('div.accessorial-panel-modal')).toBe(true);
      wrapper.find('a').simulate('click');
      expect(wrapper.exists('div.accessorial-panel-modal')).toBe(false);
    });
    it('buttons are enabled', () => {
      wrapper.find('a').simulate('click');
      expect(wrapper.find('button.button-secondary').prop('disabled')).toBe(
        false,
      );
      expect(wrapper.find('button.button-primary').prop('disabled')).toBe(
        false,
      );
    });
    it('clicking save & add another calls submitForm', () => {
      wrapper.find('a').simulate('click');
      wrapper.find('button.button-secondary').simulate('click');
      expect(submit.mock.calls.length).toBe(1);
    });
    it('clicking save & close calls submitForm', () => {
      wrapper.find('a').simulate('click');
      wrapper.find('button.button-primary').simulate('click');
      expect(submit.mock.calls.length).toBe(1);
    });
  });
});

import React from 'react';
import { shallow } from 'enzyme';

import { Creator } from './Creator';
import { no_op } from 'shared/utils';

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
const submit = jest.fn();
const clear = jest.fn();
const onFormActivation = jest.fn();
let wrapper;
describe('given a Creator', () => {
  describe('when the form is disabled', () => {
    beforeEach(() => {
      wrapper = shallow(
        <Creator
          tariff400ngItems={tariff400ng_items}
          submitForm={submit}
          formEnabled={false}
          hasSubmitSucceeded={false}
          savePreApprovalRequest={no_op}
          clearForm={clear}
          onFormActivation={onFormActivation}
        />,
      );
    });

    it('renders without crashing', () => {
      // eslint-disable-next-line
      expect(wrapper.exists('a')).toBe(true);
    });
    it('can toggle to show the form', () => {
      wrapper.find('a').simulate('click');
      expect(wrapper.exists('div.pre-approval-panel-modal')).toBe(true);
    });
    it('cancel closes the form', () => {
      wrapper.find('a').simulate('click');
      expect(wrapper.exists('div.pre-approval-panel-modal')).toBe(true);
      wrapper.find('a').simulate('click');
      expect(wrapper.exists('div.pre-approval-panel-modal')).toBe(false);
    });
    it('buttons are disabled', () => {
      wrapper.find('a').simulate('click');
      expect(wrapper.find('button.usa-button-primary').prop('disabled')).toBeTruthy();
      expect(wrapper.find('button.usa-button-secondary').prop('disabled')).toBeTruthy();
    });
  });
  describe('when the form is enabled', () => {
    beforeEach(() => {
      wrapper = shallow(
        <Creator
          tariff400ngItems={tariff400ng_items}
          submitForm={submit}
          formEnabled={true}
          hasSubmitSucceeded={false}
          savePreApprovalRequest={no_op}
          clearForm={clear}
          onFormActivation={onFormActivation}
        />,
      );
    });

    it('renders without crashing', () => {
      // eslint-disable-next-line
      expect(wrapper.exists('a')).toBe(true);
    });
    it('can toggle to show the form', () => {
      wrapper.find('a').simulate('click');
      expect(wrapper.exists('div.pre-approval-panel-modal')).toBe(true);
      expect(wrapper.state().showForm).toBe(true);
    });

    describe('when the form is open', () => {
      beforeEach(() => {
        submit.mockClear();
        wrapper.setState({ showForm: true });
      });
      it('cancel closes the form', () => {
        wrapper.find('a').simulate('click');
        expect(wrapper.exists('div.pre-approval-panel-modal')).toBe(false);
      });
      it('buttons are enabled', () => {
        expect(wrapper.find('button.usa-button-secondary').prop('disabled')).toBe(false);
        expect(wrapper.find('button.usa-button-primary').prop('disabled')).toBe(false);
      });
      it('clicking save & add another calls submitForm', () => {
        wrapper.find('button.usa-button-secondary').simulate('click');
        expect(submit.mock.calls.length).toBe(1);
      });
      it('clicking save & close calls submitForm', () => {
        wrapper.find('button.usa-button-primary').simulate('click');
        expect(submit.mock.calls.length).toBe(1);
      });
      it('after submission of click and add another, the form is cleared', () => {
        wrapper.find('button.usa-button-secondary').simulate('click');
        //redux-form will be sending this prop
        wrapper.setProps({ hasSubmitSucceeded: true }, () => {
          expect(clear.mock.calls.length).toBe(1);
          expect(wrapper.state().showForm).toBe(true);
        });
      });
      it('after submission of click and close, the form is closed', () => {
        wrapper.find('button.usa-button-primary').simulate('click');
        //redux-form will be sending this prop
        wrapper.setProps({ hasSubmitSucceeded: true }, () => {
          expect(wrapper.state().showForm).toBe(false);
        });
      });
    });
  });
});

import { shallow } from 'enzyme';
import React from 'react';
import CustomerAgreement from './CustomerAgreement';

describe('Customer Agreement', () => {
  let wrapper, props, onChange, agreementText;
  beforeEach(() => {
    onChange = jest.fn();
    agreementText = 'Text';
    props = { onChange: onChange, checked: false, agreementText: agreementText };
    wrapper = shallow(<CustomerAgreement {...props} />);
    window.alert = jest.fn();
  });

  afterEach(() => {
    window.alert = alert;
  });

  describe('checkbox', () => {
    it('is not checked on initial render', () => {
      const checkbox = wrapper.find({ type: 'checkbox' });

      expect(checkbox.props().checked).toBe(false);
    });

    it('is checked when checked prop is true', () => {
      const wrapper = shallow(<CustomerAgreement {...props} checked={true} />);
      const checkbox = wrapper.find({ type: 'checkbox' });

      expect(checkbox.props().checked).toBe(true);
    });

    it('calls onChangeHandler on change events', function() {
      const onChangeHandler = jest.fn();
      const wrapper = shallow(<CustomerAgreement {...props} onChange={onChangeHandler} />);
      const checkbox = wrapper.find({ type: 'checkbox' });
      checkbox.simulate('change', { target: { checked: true } });

      expect(onChangeHandler).toHaveBeenCalledWith(true);
    });
  });
  describe('legal agreement', () => {
    it('click triggers alert with agreementText', () => {
      wrapper.find('a').simulate('click', { preventDefault: jest.fn() });

      expect(window.alert).toHaveBeenCalledWith(props.agreementText);
    });
  });
});

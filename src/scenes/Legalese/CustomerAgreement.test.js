import { shallow } from 'enzyme';
import React from 'react';
import CustomerAgreement from './CustomerAgreement';
import CheckBox from 'shared/CheckBox';
import PopUp from 'shared/PopUp';

describe('Customer Agreement', () => {
  let wrapper, props, onAcceptTermsChange, agreementText;
  beforeEach(() => {
    onAcceptTermsChange = jest.fn();
    agreementText = 'Text';
    props = { onAcceptTermsChange: onAcceptTermsChange, checked: false, agreementText: agreementText };
    wrapper = shallow(<CustomerAgreement {...props} />);
  });

  describe('CheckBox', () => {
    it('is not checked on initial render', () => {
      const checkbox = wrapper.find(CheckBox);

      expect(checkbox.props().checked).toBe(false);
    });

    it('handleAcceptTermsChange calls props.onAcceptTermsChange', () => {
      wrapper.instance().handleAcceptTermsChange(true);

      expect(onAcceptTermsChange).toHaveBeenCalledWith(true);
    });
  });
  describe('PopUp', () => {
    it('alertMessage uses agreement text', () => {
      const popup = wrapper.find(PopUp);

      expect(popup.props().alertMessage).toBe(agreementText);
    });
  });
});

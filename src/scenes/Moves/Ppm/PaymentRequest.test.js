import React from 'react';
import { PaymentRequest } from './PaymentRequest';
import { shallow } from 'enzyme';
import CustomerAgreement from 'scenes/Legalese/CustomerAgreement';

describe('Payment Request', () => {
  let wrapper, instance;
  beforeEach(() => {
    const props = {
      currentPpm: { status: 'APPROVED', id: '1' },
      docTypes: [],
      moveDocuments: [{ id: 1 }],
      genericMoveDocSchema: {},
      moveDocSchema: {},
      updatingPPM: false,
      updateError: false,
      match: { params: { moveId: 1 } },
      getMoveDocumentsForMove: jest.fn(),
      submitExpenseDocs: jest.fn(),
    };
    wrapper = shallow(<PaymentRequest {...props} />);
    instance = wrapper.instance();
  });
  describe('PPM is approved, updatingPPM is false, and has uploaded at least one move doc', () => {
    describe('CustomerAgreement has not been accepted', () => {
      it('checkbox is not checked', () => {
        const checkbox = wrapper.find(CustomerAgreement);

        expect(checkbox.props().checked).toBe(false);
      });
      it('submit button is disabled', () => {
        wrapper.setState({ acceptTerms: true });
        const button = wrapper.find('button');

        expect(button.props().disabled).toBe(false);
      });
    });
    describe('CustomerAgreement has been accepted', () => {
      it('submit button is active', () => {
        wrapper.setState({ acceptTerms: true });
        const button = wrapper.find('button');
        const checkbox = wrapper.find(CustomerAgreement);

        expect(button.props().disabled).toBe(false);
        expect(checkbox.props().checked).toBe(true);
      });
      it('click calls applyClickHandlers', () => {
        instance.applyClickHandlers = jest.fn();
        wrapper.setState({ acceptTerms: true });
        wrapper.find('button').simulate('click');

        expect(instance.applyClickHandlers).toHaveBeenCalled();
      });
      it('applyClickHandlers calls submitDocs, submitCertificate', () => {
        instance.submitDocs = jest.fn(() => Promise.resolve());
        instance.submitCertificate = jest.fn(() => Promise.resolve());
        instance.applyClickHandlers();

        expect(instance.submitCertificate).toHaveBeenCalled();
        expect(instance.submitDocs).toHaveBeenCalled();
      });
    });
  });
  describe('No move docs uploaded', () => {
    let wrapper;
    beforeEach(() => {
      const props = {
        currentPpm: { status: 'APPROVED', id: '1' },
        docTypes: [],
        moveDocuments: [],
        genericMoveDocSchema: {},
        moveDocSchema: {},
        updatingPPM: false,
        updateError: false,
        match: { params: { moveId: 1 } },
        getMoveDocumentsForMove: jest.fn(),
        submitExpenseDocs: jest.fn(),
      };
      wrapper = shallow(<PaymentRequest {...props} />);
    });
    it('submit button is disabled even if CustomerAgreement has been accepted', () => {
      wrapper.setState({ acceptTerms: true });
      const button = wrapper.find('button');

      expect(button.props().disabled).toBe(true);
    });
  });
});

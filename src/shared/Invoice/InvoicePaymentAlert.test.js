import React from 'react';
import { shallow } from 'enzyme';
import InvoicePaymentAlert from './InvoicePaymentAlert';
import { isError, isLoading, isSuccess } from 'shared/constants';

describe('Invoice Payment Component tests', () => {
  let wrapper;

  describe('When invoice status is in processing', () => {
    it('renders under processing view ', () => {
      wrapper = shallow(<InvoicePaymentAlert createInvoiceStatus={isLoading} />);
      expect(wrapper.find('.warning--header').text()).toEqual('Sending information to USBank/Syncada.');
    });
  });
  describe('When invoice status is in failed condition', () => {
    it('renders under invoice failed view ', () => {
      wrapper = shallow(<InvoicePaymentAlert createInvoiceStatus={isError} />);
      expect(wrapper.find('.warning--header').text()).toEqual('Please try again.');
    });
  });
  describe('When invoice status is approved', () => {
    it('renders under invoice approved view ', () => {
      wrapper = shallow(<InvoicePaymentAlert createInvoiceStatus={isSuccess} />);
      expect(wrapper.find('.warning--header').text()).toEqual('The invoice has been created and will be paid soon.');
    });
  });
});

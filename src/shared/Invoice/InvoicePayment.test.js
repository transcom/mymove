import React from 'react';
import { shallow } from 'enzyme';
import InvoicePayment from './InvoicePayment';

describe('Invoice Payment Component tests', () => {
  let wrapper;
  let confirm = () => {};
  let cancel = () => {};

  describe('When invoice status is in processing', () => {
    it('renders under processing view ', () => {
      wrapper = shallow(
        <InvoicePayment
          approvePayment={confirm}
          cancelPayment={cancel}
          allowPayment={true}
          createInvoiceStatus={{
            error: null,
            isLoading: true,
            isSuccess: false,
          }}
        />,
      );
      expect(wrapper.find('.warning--header').text()).toEqual('Sending information to USBank/Syncada.');
    });
  });
  // describe('When invoice status is in confirmation', () => {
  //   it('renders under confirmation view ', () => {
  //     wrapper = shallow(
  //       <InvoicePayment
  //         approvePayment={confirm}
  //         cancelPayment={cancel}
  //         createInvoiceStatus={{
  //           error: null,
  //           isLoading: false,
  //           isSuccess: false,
  //         }}
  //       />,
  //     );
  //     expect(wrapper.find('.warning--header').text()).toEqual("Please make sure you've double-checked everything.");
  //     expect(wrapper.find('.usa-button-secondary').text()).toEqual('Cancel');
  //     expect(wrapper.find('.usa-button-primary').text()).toEqual('Approve');
  //   });
  // });
  describe('When invoice status is in failed condition', () => {
    it('renders under invoice failed view ', () => {
      wrapper = shallow(
        <InvoicePayment
          approvePayment={confirm}
          cancelPayment={cancel}
          allowPayment={true}
          createInvoiceStatus={{
            error: 'some error',
            isLoading: false,
            isSuccess: false,
          }}
        />,
      );
      expect(wrapper.find('.warning--header').text()).toEqual('Please try again.');
    });
  });
  describe('When invoice status is approved', () => {
    it('renders under invoice approved view ', () => {
      wrapper = shallow(
        <InvoicePayment
          approvePayment={confirm}
          cancelPayment={cancel}
          allowPayment={true}
          createInvoiceStatus={{
            error: null,
            isLoading: false,
            isSuccess: true,
          }}
        />,
      );
      expect(wrapper.find('.warning--header').text()).toEqual('The invoice has been created and will be paid soon.');
    });
  });
});

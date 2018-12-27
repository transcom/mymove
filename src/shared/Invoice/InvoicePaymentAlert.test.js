import React from 'react';
import { shallow } from 'enzyme';
import InvoicePaymentAlert from './InvoicePaymentAlert';
import { get } from 'lodash';

describe('Invoice Payment Component tests', () => {
  let wrapper;

  describe('When invoice status is in processing', () => {
    it('renders under processing view ', () => {
      wrapper = shallow(
        <InvoicePaymentAlert
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
  describe('When invoice status is in failed condition', () => {
    it('renders under invoice failed view ', () => {
      wrapper = shallow(
        <InvoicePaymentAlert
          createInvoiceStatus={{
            error: 'some error',
            isLoading: false,
            isSuccess: false,
          }}
        />,
      );
      expect(wrapper.find('.warning--header').text()).toEqual('Please try again.');
    });
    describe('and the api response is 409 and invoice status is SUBMITTED', () => {
      it('renders invoice already processed by another user', () => {
        wrapper = shallow(
          <InvoicePaymentAlert
            createInvoiceStatus={{
              error: {
                response: {
                  status: 409,
                  response: {
                    body: {
                      status: 'SUBMITTED',
                      approver_first_name: 'Leo',
                      approver_last_name: 'Spaceman',
                    },
                  },
                },
              },
              isLoading: false,
              isSuccess: false,
            }}
          />,
        );
        expect(wrapper.find('.warning--header').text()).toEqual(
          'Counselor Leo Spaceman already approved this invoice.',
        );
      });
    });
    describe('and the api response is 409 and invoice status is IN_PROCESS', () => {
      it('renders invoice already processed by another user', () => {
        wrapper = shallow(
          <InvoicePaymentAlert
            createInvoiceStatus={{
              error: {
                response: {
                  status: 409,
                  response: {
                    body: {
                      status: 'IN_PROCESS',
                      approver_first_name: 'Leo',
                      approver_last_name: 'Spaceman',
                    },
                  },
                },
              },
              isLoading: false,
              isSuccess: false,
            }}
          />,
        );
        expect(wrapper.find('.warning--header').text()).toEqual(
          'Counselor Leo Spaceman already submitted this invoice. Please reload your screen to see updated information.',
        );
      });
    });
    describe('and the api response is 409 and invoice status is DRAFT', () => {
      it('renders invoice already processed by another user', () => {
        wrapper = shallow(
          <InvoicePaymentAlert
            createInvoiceStatus={{
              error: {
                response: {
                  status: 409,
                  response: {
                    body: {
                      status: 'DRAFT',
                      approver_first_name: 'Leo',
                      approver_last_name: 'Spaceman',
                    },
                  },
                },
              },
              isLoading: false,
              isSuccess: false,
            }}
          />,
        );
        expect(wrapper.find('.warning--header').text()).toEqual(
          'Counselor Leo Spaceman already submitted this invoice. Please reload your screen to see updated information.',
        );
      });
    });
  });
  describe('When invoice status is approved', () => {
    it('renders under invoice approved view ', () => {
      wrapper = shallow(
        <InvoicePaymentAlert
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

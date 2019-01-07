import React from 'react';
import { shallow } from 'enzyme';
import moment from 'moment';

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
    describe('and the api response is 409 and invoice status is SUBMITTED', () => {
      let invoiceDate = '12/12/2018 12:12:00z';
      let momentDate = moment(invoiceDate);
      it('renders invoice already processed by another user', () => {
        wrapper = shallow(
          <InvoicePaymentAlert
            createInvoiceStatus={isError}
            lastInvoiceError={{
              response: {
                status: 409,
                response: {
                  body: {
                    status: 'SUBMITTED',
                    approver_first_name: 'Leo',
                    approver_last_name: 'Spaceman',
                    invoiced_date: invoiceDate,
                  },
                },
              },
            }}
          />,
        );
        expect(wrapper.find('.warning--header').text()).toEqual(
          `Leo Spaceman approved this invoice on ${momentDate.format('DD-MMM-YYYY')} at ${momentDate.format('kk:mm')}.`,
        );
      });
    });
    describe('and the api response is 409 and invoice status is IN_PROCESS', () => {
      let invoiceDate = '12/12/2018 12:12:00z';
      let momentDate = moment(invoiceDate);
      it('renders invoice already processed by another user', () => {
        wrapper = shallow(
          <InvoicePaymentAlert
            createInvoiceStatus={isError}
            lastInvoiceError={{
              response: {
                status: 409,
                response: {
                  body: {
                    status: 'IN_PROCESS',
                    approver_first_name: 'Leo',
                    approver_last_name: 'Spaceman',
                    invoiced_date: invoiceDate,
                  },
                },
              },
            }}
          />,
        );
        expect(wrapper.find('.warning--header').text()).toEqual(
          `Leo Spaceman submitted this invoice on ${momentDate.format('DD-MMM-YYYY')} at ${momentDate.format(
            'kk:mm',
          )}.`,
        );
      });
    });
    describe('and the api response is 409 and invoice status is DRAFT', () => {
      let invoiceDate = '12/12/2018 12:12:00z';
      let momentDate = moment(invoiceDate);
      it('renders invoice already processed by another user', () => {
        wrapper = shallow(
          <InvoicePaymentAlert
            createInvoiceStatus={isError}
            lastInvoiceError={{
              response: {
                status: 409,
                response: {
                  body: {
                    status: 'DRAFT',
                    approver_first_name: 'Leo',
                    approver_last_name: 'Spaceman',
                    invoiced_date: invoiceDate,
                  },
                },
              },
            }}
          />,
        );
        expect(wrapper.find('.warning--header').text()).toEqual(
          `Leo Spaceman submitted this invoice on ${momentDate.format('DD-MMM-YYYY')} at ${momentDate.format(
            'kk:mm',
          )}.`,
        );
      });
    });
  });
  describe('When invoice status is approved', () => {
    it('renders under invoice approved view ', () => {
      wrapper = shallow(<InvoicePaymentAlert createInvoiceStatus={isSuccess} />);
      expect(wrapper.find('.warning--header').text()).toEqual('The invoice has been created and will be paid soon.');
    });
  });
});

import { screen, render } from '@testing-library/react';

import e from 'constants/MoveHistory/EventTemplates/UpdatePaymentRequest/updatePaymentRequest';
import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('when a payment request has an update', () => {
  const historyRecord = {
    action: 'UPDATE',
    tableName: 'payment_requests',
    eventName: '',
    changedValues: {
      status: 'SENT_TO_GEX',
    },
    oldValues: {
      payment_request_number: '4462-6355-3',
    },
  };

  const historyRecord2 = {
    action: 'UPDATE',
    tableName: 'payment_requests',
    eventName: '',
    changedValues: {
      status: 'TPPS_RECEIVED',
    },
    oldValues: {
      payment_request_number: '4462-6355-3',
    },
  };

  const historyRecord3 = {
    action: 'UPDATE',
    tableName: 'payment_requests',
    eventName: '',
    changedValues: {
      status: 'PAID',
    },
    oldValues: {
      payment_request_number: '4462-6355-3',
    },
  };

  const historyRecordWithError = {
    action: 'UPDATE',
    tableName: 'payment_requests',
    eventName: '',
    changedValues: {
      status: 'EDI_ERROR',
    },
    oldValues: {
      payment_request_number: '4462-6355-3',
    },
  };

  it('should match the given event to the proper template', () => {
    const template = getTemplate(historyRecord);

    expect(template).toMatchObject(e);
  });

  it('should display the proper event name with correct payment request number', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    expect(screen.getByText('Updated payment request 4462-6355-3')).toBeInTheDocument();
  });

  describe('should display the proper labeled details when payment status is changed', () => {
    it.each([
      ['Status', ': Sent to GEX', historyRecord],
      ['Status', ': TPPS Received', historyRecord2],
      ['Status', ': TPPS Paid', historyRecord3],
      ['Status', ': EDI Error', historyRecordWithError],
    ])('label `%s` should have value `%s`', (label, value, record) => {
      const template = getTemplate(record);
      render(template.getDetails(record));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

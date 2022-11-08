import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdatePaymentRequestStatus/updatePaymentRequestStatus';

describe('When given a updatePaymentRequestStatus event', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updatePaymentRequestStatus',
    tableName: 'payment_requests',
    oldValues: {
      payment_request_number: '4462-6355-1',
    },
    context: [
      {
        status: 'APPROVED',
        name: 'Move management',
        price: '45985',
      },
    ],
  };

  it('should match the event to the proper template', () => {
    const template = getTemplate(historyRecord);

    expect(template).toMatchObject(e);
  });

  it('should display the event name with the correct payment request number', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    expect(screen.getByText('Submitted payment request 4462-6355-1')).toBeInTheDocument();
  });

  it('should render the expected svg icons', () => {
    const template = getTemplate(historyRecord);
    const { container } = render(template.getDetails(historyRecord));

    expect(container.querySelector("[class='svg-inline--fa fa-check successCheck']")).toBeInTheDocument();
    expect(container.querySelector("[class='svg-inline--fa fa-xmark rejectTimes']")).toBeInTheDocument();
  });

  describe('For the given history record', () => {
    it.each([
      ['Approved service items total:', '$459.85'],
      ['Move management', '$0.00'],
      ['Rejected service items total:', '$0.00'],
    ])('expect `%s` to have the value `%s`', async (label, value) => {
      const template = getTemplate(historyRecord);

      render(template.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

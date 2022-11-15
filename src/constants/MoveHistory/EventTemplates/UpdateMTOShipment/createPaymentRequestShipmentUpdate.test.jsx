import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateMTOShipment/createPaymentRequestShipmentUpdate';

describe('when given a payment request is created through shipment update', () => {
  const historyRecord = {
    action: 'INSERT',
    eventName: 'updateMTOShipment',
    tableName: 'payment_requests',
    changedValues: {
      payment_request_number: '124acb123',
    },
  };
  it('correctly matches the Request shipment reweigh event', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay(historyRecord)).toEqual('Created payment request 124acb123');
  });
  it('correctly displays the details component for shipment reweigh event', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': Pending')).toBeInTheDocument();
  });
});

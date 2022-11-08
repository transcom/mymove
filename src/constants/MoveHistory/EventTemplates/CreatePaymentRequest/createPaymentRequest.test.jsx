import { screen, render } from '@testing-library/react';

import e from 'constants/MoveHistory/EventTemplates/CreatePaymentRequest/createPaymentRequest';
import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('when given a payment request is created through shipment update', () => {
  const historyRecord = {
    action: 'INSERT',
    eventName: 'createPaymentRequest',
    tableName: 'payment_requests',
    context: [
      {
        name: 'Test service',
        price: '10123',
        status: 'REQUESTED',
        shipment_id: '123',
        shipment_type: 'HHG',
        shipment_id_abbr: 'acf7b',
      },
      {
        name: 'Domestic uncrating',
        price: '5555',
        status: 'REQUESTED',
        shipment_id: '123',
        shipment_type: 'HHG',
        shipment_id_abbr: 'acf7b',
      },
      { name: 'Move management', price: '1234', status: 'REQUESTED' },
    ],
    changedValues: { payment_request_number: '2052-7586-3' },
  };

  it('correctly matches the create payment request event', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper event name with the correct shipment id', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    expect(screen.getByText('Submitted payment request 2052-7586-3')).toBeInTheDocument();
  });

  it('displays the proper shipment type and ID', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('HHG shipment #ACF7B')).toBeInTheDocument();
  });

  describe('when given a specific set of details it displays the proper labeled information in the details column', () => {
    it.each([
      ['Move services', ': Move management'],
      ['Shipment services', ': Test service, Domestic uncrating'],
    ])('for label %s it displays the proper details value %s', async (label, value) => {
      const template = getTemplate(historyRecord);

      render(template.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

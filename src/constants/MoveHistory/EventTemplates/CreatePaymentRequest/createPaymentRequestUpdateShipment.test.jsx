import { screen, render } from '@testing-library/react';

import e from 'constants/MoveHistory/EventTemplates/CreatePaymentRequest/createPaymentRequestUpdateShipment';
import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('when given a payment request is created through shipment update', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'createPaymentRequest',
    tableName: 'mto_shipments',
    context: [
      {
        shipment_type: 'HHG',
        shipment_id_abbr: 'acf7b',
        shipment_locator: 'ABC123-01',
      },
    ],
    changedValues: { distance: 30 },
  };

  it('correctly matches the create payment request event', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper payment request record', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Shipping distance')).toBeInTheDocument();
    expect(screen.getByText(': 30 miles')).toBeInTheDocument();
  });
});

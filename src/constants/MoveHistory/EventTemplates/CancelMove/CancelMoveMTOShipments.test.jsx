import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/CancelMove/CancelMoveMTOShipments';

describe('when given a Move cancellation shipment history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: {
      status: 'CANCELED',
    },
    context: [
      {
        shipment_id_abbr: '95db3',
        shipment_locator: '6K9PYC-01',
        shipment_type: 'HHG',
      },
    ],
    eventName: 'cancelMove',
    tableName: 'mto_shipments',
  };

  it('correctly matches the Move cancellation shipment event to the proper template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the correct value in the details column', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('HHG shipment #6K9PYC-01'));
    expect(screen.getByText('Status'));
    expect(screen.getByText(': CANCELED'));
  });
});

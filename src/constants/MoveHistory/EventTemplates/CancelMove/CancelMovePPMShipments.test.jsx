import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/CancelMove/CancelMovePPMShipments';

describe('when given a Move cancellation PPM shipment history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: {
      status: 'CANCELED',
    },
    context: [
      {
        shipment_id_abbr: '87db6',
        shipment_locator: '6K9PYZ-01',
        shipment_type: 'PPM',
      },
    ],
    eventName: 'cancelMove',
    tableName: 'ppm_shipments',
  };

  it('correctly matches the Move cancellation PPM shipment event to the proper template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the correct value in the details column', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('PPM shipment #6K9PYZ-01'));
    expect(screen.getByText('Status'));
    expect(screen.getByText(': CANCELED'));
  });
});

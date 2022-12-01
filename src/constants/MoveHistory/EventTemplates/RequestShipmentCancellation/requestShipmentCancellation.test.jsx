import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/RequestShipmentCancellation/requestShipmentCancellation';

describe('when given a Request shipment cancellation history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: {
      status: 'DRAFT',
    },
    eventName: 'requestShipmentCancellation',
    oldValues: { shipment_type: 'HHG' },
    tableName: 'mto_shipments',
    context: [
      {
        shipment_id_abbr: 'acf7b',
        shipment_type: 'HHG',
      },
    ],
  };

  it('correctly matches the Request shipment cancellation event to the proper template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the correct value in the details column', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Requested cancellation for HHG shipment #ACF7B'));
  });
});

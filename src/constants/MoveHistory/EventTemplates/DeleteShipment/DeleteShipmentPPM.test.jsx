import { render, screen } from '@testing-library/react';

import e from 'constants/MoveHistory/EventTemplates/DeleteShipment/DeleteShipmentPPM';
import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('When a prime deletes a shipment', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'deleteMTOShipment',
    context: [{ shipment_id_abbr: '3b475', shipment_type: 'PPM', shipment_locator: 'ABC123-01' }],
    tableName: 'ppm_shipments',
  };

  it('correctly matches to the proper template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('correctly displays a message indicating a PPM shipment was deleted', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('PPM shipment #ABC123-01 deleted')).toBeInTheDocument();
  });
});

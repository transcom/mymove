import { render, screen } from '@testing-library/react';

import e from 'constants/MoveHistory/EventTemplates/DeleteShipment/deleteShipment';
import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('When a service counselor deletes a shipment', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'deleteShipment',
    context: [{ shipment_id_abbr: '3a771', shipment_type: 'HHG' }],
    tableName: 'mto_shipments',
  };

  it('correctly matches to the proper template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('correctly displays a message indicating a HHG shipment was deleted', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('HHG shipment #3A771 deleted')).toBeInTheDocument();
  });
});

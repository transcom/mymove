import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateMTOShipment/updateMTOShipmentByServiceItemStatus';

describe('when given an mto shipment update with service item status history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: { dest_sit_auth_end_date: '2024-12-22' },
    eventName: 'updateMTOServiceItemStatus',
    oldValues: { status: 'APPROVED' },
    tableName: 'mto_shipments',
    context: [
      {
        shipment_type: 'HHG',
        shipment_id_abbr: 'acf7b',
        shipment_locator: 'ABC123-01',
        name: 'Domestic origin price',
      },
    ],
  };
  it('correctly matches to the service item status template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper value in the details field', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('HHG shipment #ABC123-01', { exact: false })).toBeInTheDocument();
    expect(screen.getByText('Domestic origin price', { exact: false })).toBeInTheDocument();
  });
});

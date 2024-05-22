import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/ApproveShipmentDiversion/approveShipmentDiversion';

describe('when given an Approved shipment diversion history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: { status: 'APPROVED' },
    eventName: 'approveShipmentDiversion',
    tableName: 'mto_shipments',
    context: [
      {
        shipment_id_abbr: '2fa5c',
        shipment_type: 'HHG',
        shipment_locator: 'ABC123-01',
      },
    ],
  };

  it('correctly matches the Approved shipment event', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('correctly displays the proper details message', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('HHG shipment #ABC123-01')).toBeInTheDocument();
  });
  it('displays the proper name in the event name display column', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    expect(screen.getByText('Approved shipment')).toBeInTheDocument();
  });
});

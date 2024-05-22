import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/RequestShipmentDiversion/requestShipmentDiversion';

describe('when a shipment diversion is requested', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: { status: 'DIVERSION_REQUESTED' },
    eventName: 'requestShipmentDiversion',
    context: [
      {
        shipment_id_abbr: '2fa5c',
        shipment_type: 'HHG',
        shipment_locator: 'ABC123-01',
      },
    ],
    tableName: 'mto_shipments',
  };

  it('correctly matches to the proper template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('correctly displays the proper details message', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Requested diversion for HHG shipment #ABC123-01')).toBeInTheDocument();
  });
});

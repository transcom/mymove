import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';

const historyRecord = {
  ACKNOWLEDGE_SHIPMENT: {
    action: a.UPDATE,
    eventName: o.acknowledgeMovesAndShipments,
    tableName: t.mto_shipments,
    context: [
      {
        shipment_type: 'HHG',
        shipment_locator: 'RQ38D4-01',
      },
    ],
    changedValues: {
      prime_acknowledged_at: '2025-04-13T12:15:33.12345+00:00',
    },
  },
};

describe('When a shipment is acknowledged by the prime', () => {
  it('displays the prime acknowledged at timestamp', () => {
    const template = getTemplate(historyRecord.ACKNOWLEDGE_SHIPMENT);
    render(template.getDetails(historyRecord.ACKNOWLEDGE_SHIPMENT));
    const label = screen.getByText('Prime Acknowledged At:');
    expect(label).toBeInTheDocument();
    const dateElement = screen.getByText('2025-04-13T12:15:33.12345+00:00');
    expect(dateElement).toBeInTheDocument();
    const shipmentInfoElement = screen.getByText('HHG shipment #RQ38D4-01');
    expect(shipmentInfoElement).toBeInTheDocument();
  });
});

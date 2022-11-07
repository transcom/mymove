import { screen, render } from '@testing-library/react';

import e from 'constants/MoveHistory/EventTemplates/UpdateMTOShipmentStatus/updateMTOShipmentStatus';
import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('when Prime user cancels a shipment', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updateMTOShipmentStatus',
    tableName: 'mto_shipments',
    context: [
      {
        shipment_id_abbr: 'acf7b',
        shipment_type: 'HHG',
      },
    ],
    changedValues: { status: 'CANCELED' },
  };

  it('matches the given history record to the expected template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper header for the given history record', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('HHG shipment #ACF7B')).toBeInTheDocument();
  });

  describe('When given a specific set of details for a cancelled shipment', () => {
    it.each([['Status', ': CANCELED']])('displays the proper details value for %s', async (label, value) => {
      const template = getTemplate(historyRecord);
      render(template.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

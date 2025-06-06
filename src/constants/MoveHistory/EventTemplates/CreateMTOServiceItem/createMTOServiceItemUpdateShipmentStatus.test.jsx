import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/CreateMTOServiceItem/createMTOServiceItemUpdateShipmentStatus';

describe('When given a shipment updated by create service item', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: {
      status: 'APPROVALS REQUESTED',
    },
    context: [
      {
        shipment_type: 'HHG',
        shipment_locator: 'RQ38D4-01',
        shipment_id_abbr: 'a1b2c',
      },
    ],
    eventName: 'createMTOServiceItem',
    tableName: 'mto_shipments',
  };

  it('correctly matches the Create basic service item event', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay(historyRecord)).toEqual('Updated shipment');
  });
  describe('when given a specific set of details', () => {
    it.each([['Status', ': APPROVALS REQUESTED']])(
      'displays the correct details value for %s',
      async (label, value) => {
        const result = getTemplate(historyRecord);
        render(result.getDetails(historyRecord));
        expect(screen.getByText(label)).toBeInTheDocument();
        expect(screen.getByText(value)).toBeInTheDocument();
      },
    );
    it('displays the correct label for shipment', () => {
      const result = getTemplate(historyRecord);
      render(result.getDetails(historyRecord));
      expect(screen.getByText('HHG shipment #RQ38D4-01')).toBeInTheDocument();
    });
  });
});

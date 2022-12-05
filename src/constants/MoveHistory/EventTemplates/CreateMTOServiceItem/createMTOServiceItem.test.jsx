import { render, screen } from '@testing-library/react';

import e from 'constants/MoveHistory/EventTemplates/CreateMTOServiceItem/createMTOServiceItem';
import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('when given a Create basic service item history record', () => {
  const historyRecord = {
    action: 'INSERT',
    changedValues: {
      reason: 'Test',
      status: 'SUBMITTED',
    },
    context: [
      {
        name: 'Domestic uncrating',
        shipment_type: 'HHG',
        shipment_id_abbr: 'a1b2c',
      },
    ],
    eventName: 'createMTOServiceItem',
    tableName: 'mto_service_items',
  };

  const template = getTemplate(historyRecord);
  it('correctly matches the create service item event', () => {
    expect(template).toMatchObject(e);
    expect(template.getEventNameDisplay()).toEqual('Requested service item');
  });
  describe('when given a specific set of details', () => {
    it.each([
      ['service_item_name', 'Domestic uncrating'],
      ['shipment_type', 'HHG'],
      ['shipment_id_display', 'A1B2C'],
      ['reason', 'Test'],
      ['status', 'SUBMITTED'],
    ])('for label %s it displays the proper details value %s', async (label, value) => {
      render(template.getDetails(historyRecord));
      expect(screen.getByText(value, { exact: false })).toBeInTheDocument();
    });
  });
});

import { render, screen } from '@testing-library/react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import e from 'constants/MoveHistory/EventTemplates/CreateMTOServiceItem/createMTOServiceItem';
import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';

describe('when given a Create basic service item history record', () => {
  const historyRecord = {
    action: a.INSERT,
    changedValues: {
      reason: 'Test',
      status: 'SUBMITTED',
    },
    context: [
      {
        name: 'Domestic uncrating',
        shipment_type: 'HHG',
      },
    ],
    eventName: o.createMTOServiceItem,
    tableName: t.mto_service_items,
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
      ['reason', 'Test'],
      ['status', 'SUBMITTED'],
    ])('for label %s it displays the proper details value %s', async (label, value) => {
      render(template.getDetails(historyRecord));
      expect(screen.getByText(value, { exact: false })).toBeInTheDocument();
    });
  });
});

import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

describe('when given a Update sit entry date service item history record', () => {
  const historyRecord = {
    action: a.UPDATE,
    changedValues: {
      sit_entry_date: '2023-10-31',
    },
    context: [
      {
        name: "Domestic destination add'l SIT",
        shipment_type: 'HHG',
        shipment_id_abbr: 'a1b2c',
        shipment_locator: 'ABC123-01',
      },
    ],
    oldValues: {
      sit_entry_date: '2023-10-01',
    },
    eventName: o.updateServiceItemSitEntryDate,
    tableName: t.mto_service_items,
  };
  it('displays shipment type, shipment ID, and service item name properly', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText("HHG shipment #ABC123-01, Domestic destination add'l SIT")).toBeInTheDocument();
  });

  describe('When given a specific set of details', () => {
    it.each([['SIT entry date', ': 31 Oct 2023']])('displays the proper details value for %s', async (label, value) => {
      const result = getTemplate(historyRecord);
      render(result.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

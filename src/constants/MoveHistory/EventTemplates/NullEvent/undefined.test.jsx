import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/NullEvent/undefined';

describe('when given an unidentifiable move history record', () => {
  describe('depending on the table name', () => {
    it.each([
      ['Updated shipment', 'mto_shipments'],
      ['Updated shipment', 'mto_agents'],
      ['Updated shipment', 'addresses'],
      ['Updated move', 'moves'],
      ['Updated payment request', 'payment_requests'],
      ['Updated allowances', 'entitlements'],
      ['Updated service item', 'mto_service_items'],
      ['Updated order', 'orders'],
    ])('it displays the `%s` event name for the table `%s`', async (eventDisplayName, table) => {
      const historyRecord = {
        action: null,
        eventName: 'testname',
        tableName: `${table}`,
      };

      const template = getTemplate(historyRecord);

      expect(template).toEqual(e);
      render(template.getEventNameDisplay(historyRecord));
      expect(screen.getByText(eventDisplayName)).toBeInTheDocument();
      render(template.getDetails(historyRecord));
      expect(screen.getByText('-')).toBeInTheDocument();
    });
  });
});

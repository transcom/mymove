import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/ReviewShipmentAddressUpdate/reviewShipmentAddressUpdateServiceItem';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

describe('when given a updated shipment address request, update move history record', () => {
  const historyRecord = {
    action: a.UPDATE,
    eventName: o.reviewShipmentAddressUpdate,
    tableName: t.mto_service_items,
    changedValues: { pricing_estimate: 12345 },
  };

  it('correctly matches the update service item status, update move event', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  describe('When given an updated shipment address request, update service item history record', () => {
    it.each([['Estimated price', ': $123.45']])('displays the proper details value for %s', async (label, value) => {
      const template = getTemplate(historyRecord);
      render(template.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

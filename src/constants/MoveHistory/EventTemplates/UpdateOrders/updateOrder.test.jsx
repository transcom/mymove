import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateOrders/updateOrder';

describe('when given an Order update history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updateOrder',
    tableName: 'orders',
    changedValues: { old_duty_location_id: 'ID1', new_duty_location_id: 'ID2', has_dependents: 'false' },
    context: [{ old_duty_location_name: 'old name', new_duty_location_name: 'new name' }],
  };
  it('correctly matches to the proper template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });
  describe('When given a specific set of details for updated orders', () => {
    it.each([
      ['New duty location name', ': new name'],
      ['Dependents included', ': false'],
    ])('displays the proper details value for %s', async (label, value) => {
      const template = getTemplate(historyRecord);
      render(template.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

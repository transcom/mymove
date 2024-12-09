import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateOrders/updateOrderUpdateAllowances';

describe('when given a update order, update allowance history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'counselingUpdateOrder',
    tableName: 'entitlements',
    changedValues: { authorized_weight: 1650 },
  };

  it('correctly matches the update order, update allowance event', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper update allowances record', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Max Billable Weight')).toBeInTheDocument();
    expect(screen.getByText(': 1,650 lbs')).toBeInTheDocument();
  });
});

import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateAllowances/updateAllowanceUpdateOrder';

describe('when given a update allowance, update order history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updateOrder',
    tableName: 'entitlements',
    changedValues: { authorized_weight: 11000 },
  };

  it('correctly matches the update allowance, update order event', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper update order record', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Authorized weight')).toBeInTheDocument();
    expect(screen.getByText(': 11,000 lbs')).toBeInTheDocument();
  });
});

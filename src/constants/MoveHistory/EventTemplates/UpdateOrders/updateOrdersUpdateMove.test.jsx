import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateOrders/updateOrdersUpdateMove';

describe('when given an mto shipment update for the moves table', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: { status: 'APPROVED' },
    eventName: 'updateOrders',
    oldValues: { status: 'APPROVALS REQUESTED' },
    tableName: 'moves',
  };
  it('correctly matches to the service item status template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper value in the details field', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': APPROVED')).toBeInTheDocument();
  });
});

import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/ApproveShipment/approveShipmentApproveMove';

describe('when given an Approved shipment, Approved move history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: { status: 'APPROVED' },
    eventName: 'approveShipment',
    oldValues: { status: 'APPROVALS REQUESTED' },
    tableName: 'moves',
  };
  it('correctly matches to the Approved shipment, Approved move template', () => {
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

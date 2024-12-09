import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/ApproveShipment/approveShipmentUpdateAllowances';

describe('when given an Approved shipment, Updated allowances history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: { authorized_weight: 13230 },
    eventName: 'approveShipment',
    tableName: 'entitlements',
  };
  it('correctly matches to the Approved shipment, Updated allowances template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper value in the details field', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Max Billable Weight')).toBeInTheDocument();
    expect(screen.getByText(': 13,230 lbs')).toBeInTheDocument();
  });
});

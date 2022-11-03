import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/ApproveShipment/approveShipment';

describe('when given an Approved shipment history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: { status: 'APPROVED' },
    eventName: 'approveShipment',
    oldValues: { shipment_type: 'HHG' },
    tableName: 'mto_shipments',
    context: [
      {
        shipment_id_abbr: '2fa5c',
        shipment_type: 'HHG',
      },
    ],
  };
  it('correctly matches to the Approved shipment template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper value in the details field', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('HHG shipment #2FA5C')).toBeInTheDocument();
  });
});

import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/ApproveShipmentDiversion/approveShipmentDiversion';

describe('when given an Approved shipment diversion history record', () => {
  const historyRecord = {
    changedValues: { status: 'APPROVED' },
    eventName: 'approveShipmentDiversion',
    context: [
      {
        shipment_id_abbr: '2fa5c',
        shipment_type: 'HHG',
      },
    ],
  };

  it('correctly matches the Approved shipment event', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('correctly displays the proper details message', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('HHG shipment #2FA5C')).toBeInTheDocument();
  });
});

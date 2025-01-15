import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

describe('When approving/rejecting a payment service item', () => {
  const rejectPaymentServiceItemRecord = {
    action: a.UPDATE,
    actionTstampClk: '2025-01-10T19:44:31.255Z',
    actionTstampStm: '2025-01-10T19:44:31.253Z',
    actionTstampTx: '2025-01-10T19:44:31.220Z',
    context: [
      {
        shipment_type: 'PPM',
        shipment_locator: 'RQ38D4-01',
        shipment_id_abbr: 'f10be',
      },
    ],
    changedValues: {
      reason: 'Some reason',
      status: 'DENIED',
    },
    eventName: o.updatePaymentServiceItemStatus,
    tableName: t.payment_service_items,
    id: '2419f1db-3f8b-4823-974f-9aa4edb753da',
    objectId: 'eee30fb1-dc66-4821-a17c-2ecf431ceb9d',
    oldValues: {
      id: 'eee30fb1-dc66-4821-a17c-2ecf431ceb9d',
      ppm_shipment_id: '86329c14-564b-4580-94b9-8a2e80bccefc',
      reason: null,
      status: null,
    },
  };

  const approvePaymentServiceItemRecord = { ...rejectPaymentServiceItemRecord };
  approvePaymentServiceItemRecord.changedValues = {
    status: 'APPROVED',
  };

  it('displays an approved payment service item', () => {
    const template = getTemplate(approvePaymentServiceItemRecord);

    render(template.getDetails(approvePaymentServiceItemRecord));
    expect(screen.getByText('PPM shipment #RQ38D4-01')).toBeInTheDocument();
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': APPROVED')).toBeInTheDocument();
  });

  it('displays a rejected payment service item and the rejection reason', () => {
    const template = getTemplate(rejectPaymentServiceItemRecord);

    render(template.getDetails(rejectPaymentServiceItemRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': DENIED')).toBeInTheDocument();
    expect(screen.getByText('Reason')).toBeInTheDocument();
    expect(screen.getByText(': Some reason')).toBeInTheDocument();
  });
});

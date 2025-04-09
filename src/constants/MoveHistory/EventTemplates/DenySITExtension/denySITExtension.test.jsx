import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import denySITExtension from 'constants/MoveHistory/EventTemplates/DenySITExtension/denySITExtension';

describe('when given a Deny SIT Extension item history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    actionTstampClk: '2025-04-09T13:43:45.092Z',
    actionTstampStm: '2025-04-09T13:43:45.092Z',
    actionTstampTx: '2025-04-09T13:43:45.081Z',
    changedValues: {
      customer_expense: false,
      decision_date: '2025-04-09T13:43:45.090591',
      office_remarks: 'rejected',
      status: 'DENIED',
    },
    context: [
      {
        shipment_id_abbr: '3d95a',
        shipment_locator: 'SITEXT-01',
        shipment_type: 'HHG',
      },
    ],
    eventName: 'denySITExtension',
    id: 'a5f69178-dd8d-4392-b087-b1325a985531',
    objectId: 'f77846bc-b5b4-47e8-9883-870806d34d4c',
    oldValues: {
      approved_days: null,
      contractor_remarks: 'remarks',
      customer_expense: null,
      decision_date: null,
      id: 'f77846bc-b5b4-47e8-9883-870806d34d4c',
      mto_shipment_id: '3d95a553-803a-4441-b8c1-48b32e9e539e',
      office_remarks: null,
      request_reason: 'SERIOUS_ILLNESS_MEMBER',
      requested_days: 10,
      status: 'PENDING',
    },
    relId: 19955,
    schemaName: 'public',
    sessionUserEmail: 'multi-role-20250409133916-e8d4c95f3042@example.com',
    sessionUserFirstName: 'Alice',
    sessionUserId: '73a87ab6-bee0-4c6f-b008-7bd1682c6e9f',
    sessionUserLastName: 'Bob',
    sessionUserTelephone: '333-333-3333',
    tableName: 'sit_extensions',
    transactionId: 2031,
  };

  it('correctly matches to the Deny SIT extension template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(denySITExtension);
  });

  it('returns the correct event display name', () => {
    expect(denySITExtension.getEventNameDisplay()).toEqual('SIT extension denied');
  });

  it('renders the Deny SIT extension details correctly', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));

    expect(screen.getByText('HHG shipment #SITEXT-01')).toBeInTheDocument();
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': DENIED')).toBeInTheDocument();
    expect(screen.getByText('Office remarks')).toBeInTheDocument();
    expect(screen.getByText(': rejected')).toBeInTheDocument();
    expect(screen.getByText('Converted to customer expense')).toBeInTheDocument();
    expect(screen.getByText(': No')).toBeInTheDocument();
    expect(screen.getByText('Decision date')).toBeInTheDocument();
    expect(screen.getByText(': 09 Apr 2025')).toBeInTheDocument();
  });
});

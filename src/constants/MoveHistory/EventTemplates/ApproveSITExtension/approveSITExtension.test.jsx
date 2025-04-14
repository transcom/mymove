import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import approveSITExtension from 'constants/MoveHistory/EventTemplates/ApproveSITExtension/approveSITExtension';
import { SIT_EXTENSION_STATUS } from 'constants/sitExtensions';
import Actions from 'constants/MoveHistory/Database/Actions';

describe('when given an Approve SIT Extension item history record', () => {
  const historyRecord = {
    action: Actions.UPDATE,
    changedValues: {
      decision_date: '2025-04-09T13:43:45.090591',
      office_remarks: 'approved',
      status: SIT_EXTENSION_STATUS.APPROVED,
      approved_days: 10,
    },
    context: [
      {
        shipment_id_abbr: '3d95a',
        shipment_locator: 'SITEXT-01',
        shipment_type: 'HHG',
      },
    ],
    eventName: 'approveSITExtension',
    tableName: 'sit_extensions',
  };

  it('correctly matches to the Approve SIT extension template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(approveSITExtension);
  });

  it('returns the correct event display name', () => {
    expect(approveSITExtension.getEventNameDisplay()).toEqual('SIT extension approved');
  });

  it('renders the Approve SIT extension details correctly', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));

    expect(screen.getByText('HHG shipment #SITEXT-01')).toBeInTheDocument();
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': APPROVED')).toBeInTheDocument();
    expect(screen.getByText('Office remarks')).toBeInTheDocument();
    expect(screen.getByText(': approved')).toBeInTheDocument();
    expect(screen.getByText('Approved days')).toBeInTheDocument();
    expect(screen.getByText(': 10')).toBeInTheDocument();
    expect(screen.getByText('Decision date')).toBeInTheDocument();
    expect(screen.getByText(': 09 Apr 2025')).toBeInTheDocument();
  });
});

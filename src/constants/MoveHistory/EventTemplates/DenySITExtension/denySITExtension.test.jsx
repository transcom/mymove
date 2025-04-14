import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import denySITExtension from 'constants/MoveHistory/EventTemplates/DenySITExtension/denySITExtension';
import Actions from 'constants/MoveHistory/Database/Actions';
import { SIT_EXTENSION_STATUS } from 'constants/sitExtensions';

describe('when given a Deny SIT Extension item history record', () => {
  const historyRecord = {
    action: Actions.UPDATE,
    changedValues: {
      customer_expense: false,
      decision_date: '2025-04-09T13:43:45.090591',
      office_remarks: 'rejected',
      status: SIT_EXTENSION_STATUS.DENIED,
    },
    context: [
      {
        shipment_id_abbr: '3d95a',
        shipment_locator: 'SITEXT-01',
        shipment_type: 'HHG',
      },
    ],
    eventName: 'denySITExtension',
    tableName: 'sit_extensions',
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

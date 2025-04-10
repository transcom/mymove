import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import createSITExtension from 'constants/MoveHistory/EventTemplates/CreateSITExtension/createSITExtension';

describe('when given a Create SIT Extension item history record', () => {
  const historyRecord = {
    action: 'INSERT',
    changedValues: {
      contractor_remarks: 'remarks',
      request_reason: 'SERIOUS_ILLNESS_MEMBER',
      requested_days: 10,
      status: 'PENDING',
    },
    context: [
      {
        shipment_id_abbr: '3d95a',
        shipment_locator: 'SITEXT-01',
        shipment_type: 'HHG',
      },
    ],
    eventName: 'createSITExtension',
    tableName: 'sit_extensions',
  };

  it('correctly matches to the Create SIT extension template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(createSITExtension);
  });

  it('returns the correct event display name', () => {
    expect(createSITExtension.getEventNameDisplay()).toEqual('SIT extension requested');
  });

  it('renders the Create SIT extension details correctly', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));

    expect(screen.getByText('HHG shipment #SITEXT-01')).toBeInTheDocument();
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': PENDING')).toBeInTheDocument();
    expect(screen.getByText('Contractor remarks')).toBeInTheDocument();
    expect(screen.getByText(': remarks')).toBeInTheDocument();
    expect(screen.getByText('Request reason')).toBeInTheDocument();
    expect(screen.getByText(': Serious illness of the member')).toBeInTheDocument();
    expect(screen.getByText('Requested days')).toBeInTheDocument();
    expect(screen.getByText(': 10')).toBeInTheDocument();
  });
});

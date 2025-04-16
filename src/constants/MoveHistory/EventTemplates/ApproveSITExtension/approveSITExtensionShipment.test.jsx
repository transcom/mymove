import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import approveSITExtensionShipment from 'constants/MoveHistory/EventTemplates/ApproveSITExtension/approveSITExtensionShipment';
import Actions from 'constants/MoveHistory/Database/Actions';

describe('when given an Approve SIT Extension Shipment item history record', () => {
  const historyRecord = {
    action: Actions.UPDATE,
    changedValues: {
      dest_sit_auth_end_date: '2025-11-13',
    },
    context: [
      {
        shipment_id_abbr: '3d95a',
        shipment_locator: 'SITEXT-01',
        shipment_type: 'HHG',
      },
    ],
    eventName: 'approveSITExtension',
    tableName: 'mto_shipments',
  };

  it('correctly matches to the Approve SIT extension Shipment template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(approveSITExtensionShipment);
  });

  it('returns the correct event display name', () => {
    expect(approveSITExtensionShipment.getEventNameDisplay()).toEqual('SIT extension approved');
  });

  it('renders the Approve SIT extension Shipment details correctly', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));

    expect(screen.getByText('HHG shipment #SITEXT-01')).toBeInTheDocument();
    expect(screen.getByText('Destination SIT authorized date')).toBeInTheDocument();
    expect(screen.getByText(': 2025-11-13')).toBeInTheDocument();
  });
});

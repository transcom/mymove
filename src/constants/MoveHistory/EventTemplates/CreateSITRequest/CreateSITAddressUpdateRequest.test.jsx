import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

describe('When given an update Destination SIT address history record', () => {
  const historyRecord = {
    action: a.INSERT,
    changedValues: {
      contractor_remarks: 'Need to store in Florence',
      status: 'REQUESTED',
    },
    context: [
      {
        name: 'Domestic destination SIT delivery',
        shipment_type: 'HHG',
        shipment_locator: 'RQ38D4-01',
        shipment_id_abbr: 'e4285',
        sit_destination_address_final: `{"id":"14a265d6-95b4-4842-a2ed-e020ba7da3fb","street_address_1":"676 Destination Sit Req","street_address_2":null,"city":"Florence","state":"MT","postal_code":"59805","created_at":"2023-11-21T02:56:56.832038","updated_at":"2023-11-21T02:56:56.832038","street_address_3":null,"country":null}`,
        sit_destination_address_initial: `{"id":"ff666bfe-1a2c-45e0-b38a-18c138958f16","street_address_1":"4 Delivery address init","street_address_2":null,"city":"Great Falls","state":"MT","postal_code":"59402","created_at":"2023-11-21T02:56:08.299416","updated_at":"2023-11-21T02:56:08.299416","street_address_3":null,"country":null}`,
      },
    ],
    eventName: o.createSITAddressUpdateRequest,
    tableName: t.sit_address_updates,
  };

  it('displays shipment type, shipment ID, and service item name properly', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('HHG shipment #RQ38D4-01, Domestic destination SIT delivery')).toBeInTheDocument();
  });

  it('displays the status of the request', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': Updated')).toBeInTheDocument();
  });

  it('displays the initial SIT destination address', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Final SIT delivery address')).toBeInTheDocument();
    expect(screen.getByText(': 4 Delivery address init, Great Falls, MT 59402')).toBeInTheDocument();
  });

  it('displays the final SIT destination address', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Initial SIT delivery address')).toBeInTheDocument();
    expect(screen.getByText(': 676 Destination Sit Req, Florence, MT 59805')).toBeInTheDocument();
  });

  it('displays the contractor remarks', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Contractor remarks')).toBeInTheDocument();
    expect(screen.getByText(': Need to store in Florence')).toBeInTheDocument();
  });
});

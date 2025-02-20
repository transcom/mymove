import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateMTOShipment/updateMTOShipmentUpdateAllowance';

describe('when given an update to the allowance due to MTOShipment update, update MTO shipment history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updateMTOShipment',
    tableName: 'entitlements',
    context: [
      {
        shipment_type: 'HHG',
        shipment_id_abbr: 'acf7b',
        shipment_locator: 'ABC123-01',
      },
    ],
    changedValues: { authorized_weight: 1650 },
  };

  it('correctly matches the update to the allowance, update MTO shipment event', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper update MTO shipment record', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Max billable weight')).toBeInTheDocument();
    expect(screen.getByText(': 1,650 lbs')).toBeInTheDocument();
  });
});

import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/ApproveShipment/createStandardServiceItem';

describe('when given a Create standard service item history record', () => {
  const historyRecord = {
    action: 'INSERT',
    context: [
      {
        shipment_type: 'HHG',
        shipment_id_abbr: 'a1b2c',
        name: 'Domestic linehaul',
      },
    ],
    eventName: 'approveShipment',
    tableName: 'mto_service_items',
  };

  it('correctly matches the Create standard service item template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the correct values in the details column', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('HHG shipment #A1B2C, Domestic linehaul')).toBeInTheDocument();
  });
});

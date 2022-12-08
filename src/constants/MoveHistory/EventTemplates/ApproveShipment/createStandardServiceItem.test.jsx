import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';

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
  it('correctly matches the Create standard service item event', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));
    expect(template.getEventNameDisplay(template)).toEqual('Approved service item');
    expect(screen.getByText('HHG shipment #A1B2C', { exact: false })).toBeInTheDocument();
    expect(screen.getByText('Domestic linehaul', { exact: false })).toBeInTheDocument();
  });
});

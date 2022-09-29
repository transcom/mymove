import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/approveShipment';

describe('when given an Approved shipment history record', () => {
  const item = {
    action: 'UPDATE',
    changedValues: { status: 'APPROVED' },
    eventName: 'approveShipment',
    oldValues: { shipment_type: 'HHG' },
    tableName: 'mto_shipments',
  };
  it('correctly matches the Approved shipment event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    render(result.getDetails(item));
    expect(screen.getByText('HHG shipment')).toBeInTheDocument();
  });
});

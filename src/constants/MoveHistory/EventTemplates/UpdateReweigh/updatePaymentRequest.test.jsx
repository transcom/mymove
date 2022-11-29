import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('reweighs update', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updateReweigh',
    tableName: 'payment_requests',
    context: [{ shipment_type: 'HHG' }],
    changedValues: { recalculation_of_payment_request_id: '1234' },
    oldValues: { payment_request_number: '0288-7994-1' },
  };
  it('correctly matches the reweigh payment request', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': Recalculated payment request')).toBeInTheDocument();
  });
});

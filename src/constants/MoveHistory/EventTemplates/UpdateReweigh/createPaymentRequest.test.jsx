import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('when given a payment request is created through reweigh', () => {
  const historyRecord = {
    action: 'INSERT',
    eventName: 'updateReweigh',
    tableName: 'payment_requests',
  };
  it('correctly matches the Request shipment reweigh event', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': Pending')).toBeInTheDocument();
  });
});

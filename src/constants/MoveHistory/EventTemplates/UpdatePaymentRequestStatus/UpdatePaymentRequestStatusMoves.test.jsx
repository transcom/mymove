import { screen, render } from '@testing-library/react';

import e from 'constants/MoveHistory/EventTemplates/UpdatePaymentRequestStatus/UpdatePaymentRequestStatusMoves';
import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('When given a completed services counseling for a move', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updatePaymentRequestStatus',
    tableName: 'moves',
  };
  it('correctly matches the update mto status services counseling completed event to the proper template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper name in the event name display column', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    expect(screen.getByText('Updated move')).toBeInTheDocument();
  });

  it('displays default when TIO ID is not present', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Payment Requests Addressed')).toBeInTheDocument();
  });

  it('displays correct details when a TIO is unassigned', () => {
    historyRecord.changedValues = {
      ...historyRecord.changedValues,
      tio_payment_request_assigned_id: null,
    };
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Payment Requests Addressed')).toBeInTheDocument();
    expect(screen.getByText('Task invoicing officer unassigned')).toBeInTheDocument();
  });
});

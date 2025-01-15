import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/AcknowledgeExcessUnaccompaniedBaggageWeightRisk/acknowledgeExcessUnaccompaniedBaggageWeightRisk';

describe('when given an Acknowledge excess unaccompanied baggage weight risk history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'acknowledgeExcessUnaccompaniedBaggageWeightRisk',
    tableName: 'moves',
  };

  it('correctly matches the Acknowledge excess unaccompanied baggage weight risk template', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(e);
  });
  it('renders the default details in the details column when excess risk key is not present ', () => {
    const newHistoryRecord = {
      ...historyRecord,
      changedValues: { status: 'APPROVED' },
    };
    const template = getTemplate(newHistoryRecord);
    render(template.getDetails(newHistoryRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': APPROVED')).toBeInTheDocument();
  });

  it('renders the proper message in the details column when excess_unaccompanied_baggage_weight_acknowledged_at is present ', () => {
    const newHistoryRecordAcknowledged = {
      ...historyRecord,
      changedValues: { excess_unaccompanied_baggage_weight_acknowledged_at: 'this would usually be a time value' },
    };
    const template = getTemplate(newHistoryRecordAcknowledged);
    render(template.getDetails(newHistoryRecordAcknowledged));
    expect(screen.getByText('Dismissed excess unaccompanied baggage weight alert')).toBeInTheDocument();
  });
});

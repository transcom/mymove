import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/AcknowledgeExcessWeightRisk/acknowledgeExcessWeightRisk';

describe('when given an Acknowledge excess weight risk history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'acknowledgeExcessWeightRisk',
    tableName: 'moves',
    eventNameDisplay: 'Updated move',
  };

  it('correctly matches the Acknowledge excess weight risk template', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(e);
  });
  it('it correctly matches the event that acknowledges weight risk', () => {
    const result = getTemplate(historyRecord);
    expect(result.getEventNameDisplay()).toMatch(historyRecord.eventNameDisplay);
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

  it('renders the proper message in the details column when excess_weight_acknowledged_at is present ', () => {
    const newHistoryRecordAcknowledged = {
      ...historyRecord,
      changedValues: { excess_weight_acknowledged_at: 'this would usually be a time value' },
    };
    const template = getTemplate(newHistoryRecordAcknowledged);
    render(template.getDetails(newHistoryRecordAcknowledged));
    expect(screen.getByText('Dismissed excess weight alert')).toBeInTheDocument();
  });
});

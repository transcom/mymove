import { screen, render } from '@testing-library/react';

import e from 'constants/MoveHistory/EventTemplates/FinishDocumentReview/FinishDocumentReviewMoves';
import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('When given a completed services counseling for a move', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'finishDocumentReview',
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

  it('displays default when SC ID is not present', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('PPM Closeout Complete')).toBeInTheDocument();
    expect(screen.queryByText('Closeout Counselor Unassigned')).not.toBeInTheDocument();
  });

  it('displays correct details when a SC is unassigned', () => {
    historyRecord.changedValues = {
      ...historyRecord.changedValues,
      sc_closeout_assigned_id: null,
    };
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('PPM Closeout Complete')).toBeInTheDocument();
    expect(screen.getByText('Closeout Counselor Unassigned')).toBeInTheDocument();
  });
});

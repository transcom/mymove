import { screen, render } from '@testing-library/react';

import e from 'constants/MoveHistory/EventTemplates/UpdateMTOStatusServiceCounselingCompleted/updateMTOStatusServiceCounselingCompleted';
import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('When given a completed services counseling for a move', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updateMTOStatusServiceCounselingCompleted',
    tableName: 'moves',
  };
  it('correctly matches the update mto status services counseling completed event to the proper template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper message in the details column', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Counseling Completed')).toBeInTheDocument();
  });
  it('displays the proper name in the event name display column', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    expect(screen.getByText('Updated move')).toBeInTheDocument();
  });
  it('displays correct details when an SC is unassigned', () => {
    historyRecord.changedValues = {
      ...historyRecord.changedValues,
      sc_counseling_assigned_id: null,
    };
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Counseling Completed')).toBeInTheDocument();
    expect(screen.getByText('Counselor Unassigned')).toBeInTheDocument();
  });
});

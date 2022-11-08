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
});

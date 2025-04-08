import { screen, render } from '@testing-library/react';

import e from 'constants/MoveHistory/EventTemplates/UpdateAssignedOfficeUser/DeleteAssignedOfficeUser';
import getTemplate from 'constants/MoveHistory/TemplateManager';
import { MOVE_STATUSES } from 'shared/constants';

describe('When given a move that has been unassigned', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'deleteAssignedOfficeUser',
    tableName: 'moves',
    changedValues: {
      sc_assigned_id: null,
    },
  };

  it('correctly matches the template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper name in the event name display column', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    expect(screen.getByText('Updated move')).toBeInTheDocument();
  });

  describe('displays the proper details for', () => {
    it('closeout counselor', () => {
      const template = getTemplate(historyRecord);

      render(template.getDetails(historyRecord));
      expect(screen.getByText('Closeout counselor unassigned')).toBeInTheDocument();
    });
    it('services counselor', () => {
      const template = getTemplate(historyRecord);
      historyRecord.oldValues = { status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING };

      render(template.getDetails(historyRecord));
      expect(screen.getByText('Counselor unassigned')).toBeInTheDocument();
    });
    it('task ordering officer', () => {
      historyRecord.changedValues = { too_assigned_id: null };
      const template = getTemplate(historyRecord);

      render(template.getDetails(historyRecord));
      expect(screen.getByText('Task ordering officer unassigned')).toBeInTheDocument();
    });
    it('destination queue task ordering officer ', () => {
      historyRecord.changedValues = { too_destination_assigned_id: null };
      const template = getTemplate(historyRecord);

      render(template.getDetails(historyRecord));
      expect(screen.getByText('Task ordering officer unassigned')).toBeInTheDocument();
    });
    it('task invoicing officer', () => {
      historyRecord.changedValues = { tio_assigned_id: null };
      const template = getTemplate(historyRecord);

      render(template.getDetails(historyRecord));
      expect(screen.getByText('Task invoicing officer unassigned')).toBeInTheDocument();
    });
  });
});

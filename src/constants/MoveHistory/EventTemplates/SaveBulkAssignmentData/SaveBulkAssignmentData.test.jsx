import { screen, render } from '@testing-library/react';

import e from 'constants/MoveHistory/EventTemplates/SaveBulkAssignmentData/SaveBulkAssignmentData';
import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('When given a move that has been assigned', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'saveBulkAssignmentData',
    tableName: 'moves',
    changedValues: {
      sc_counseling_assigned_id: 'fb625e3c-067c-49d7-8fd9-88ef040e6137',
    },
    oldValues: {
      sc_counseling_assigned_id: null,
    },
    context: [{ assigned_office_user_last_name: 'Daniels', assigned_office_user_first_name: 'Jayden' }],
  };

  it('correctly matches the template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper name in the event name display column', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    expect(screen.getByText('Move assignment updated')).toBeInTheDocument();
  });

  describe('displays the proper details for', () => {
    it('services counselor', () => {
      const template = getTemplate(historyRecord);
      historyRecord.oldValues = {
        sc_counseling_assigned_id: null,
      };

      render(template.getDetails(historyRecord));
      expect(screen.getByText('Counselor assigned')).toBeInTheDocument();
      expect(screen.getByText(': Daniels, Jayden')).toBeInTheDocument();
    });
    it('closeout counselor', () => {
      const template = getTemplate(historyRecord);
      historyRecord.changedValues = { sc_closeout_assigned_id: 'fb625e3c-067c-49d7-8fd9-88ef040e6137' };
      historyRecord.oldValues = {
        sc_closeout_assigned_id: null,
      };

      render(template.getDetails(historyRecord));
      expect(screen.getByText('Closeout counselor assigned')).toBeInTheDocument();
      expect(screen.getByText(': Daniels, Jayden')).toBeInTheDocument();
    });
    it('task ordering officer', () => {
      historyRecord.changedValues = { too_assigned_id: 'fb625e3c-067c-49d7-8fd9-88ef040e6137' };
      historyRecord.oldValues = { too_assigned_id: null };
      historyRecord.context = [
        { assigned_office_user_last_name: 'Robinson', assigned_office_user_first_name: 'Brian' },
      ];

      const template = getTemplate(historyRecord);

      render(template.getDetails(historyRecord));
      expect(screen.getByText('Task ordering officer assigned')).toBeInTheDocument();
      expect(screen.getByText(': Robinson, Brian')).toBeInTheDocument();
    });
    it('task invoicing officer', () => {
      historyRecord.changedValues = { tio_payment_request_assigned_id: 'fb625e3c-067c-49d7-8fd9-88ef040e6137' };
      historyRecord.oldValues = { tio_payment_request_assigned_id: null };
      historyRecord.context = [{ assigned_office_user_last_name: 'Luvu', assigned_office_user_first_name: 'Frankie' }];

      const template = getTemplate(historyRecord);

      render(template.getDetails(historyRecord));
      expect(screen.getByText('Task invoicing officer assigned')).toBeInTheDocument();
      expect(screen.getByText(': Luvu, Frankie')).toBeInTheDocument();
    });
  });
});

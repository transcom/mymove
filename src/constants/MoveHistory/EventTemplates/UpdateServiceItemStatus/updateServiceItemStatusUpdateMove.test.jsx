import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateServiceItemStatus/updateServiceItemStatusUpdateMove';

describe('when given a update service item status, update move history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updateMTOServiceItemStatus',
    tableName: 'moves',
    oldValues: { status: 'APPROVALS REQUESTED' },
    changedValues: { status: 'APPROVED' },
  };

  it('correctly matches the update service item status, update move event', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  describe('When given an updated service item status, update move history record', () => {
    it.each([['Status', ': APPROVED']])('displays the proper details value for %s', async (label, value) => {
      const template = getTemplate(historyRecord);
      render(template.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateMoveTaskOrderStatus/updateMoveTaskOrderStatus';

describe('when given a Move approved history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: {
      available_to_prime_at: '2022-04-13T15:21:31.746028+00:00',
      status: 'APPROVED',
    },
    eventName: 'updateMoveTaskOrderStatus',
    tableName: 'moves',
  };

  it('correctly matches the Update move task order to the proper template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the correct value in the event name and details column', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Approved move'));
    expect(screen.getByText('Created Move Task Order (MTO)'));
  });

  it('displays correct details when a TOO is unassigned', () => {
    historyRecord.changedValues = {
      ...historyRecord.changedValues,
      too_task_order_assigned_id: null,
    };

    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Created Move Task Order (MTO)'));
    expect(screen.getByText('Task ordering officer unassigned'));
  });
  it('displays correct details when a TOO is unassigned on the destination request queue', () => {
    historyRecord.changedValues = {
      ...historyRecord.changedValues,
      too_destination_assigned_id: null,
    };
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Created Move Task Order (MTO)'));
    expect(screen.getByText('Destination task ordering officer unassigned'));
  });
});

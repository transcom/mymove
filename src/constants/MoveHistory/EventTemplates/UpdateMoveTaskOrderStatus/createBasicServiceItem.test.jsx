import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateMoveTaskOrderStatus/createBasicServiceItem';

describe('when given a Create basic service item history record', () => {
  const historyRecord = {
    action: 'INSERT',
    context: [
      {
        name: 'Move management',
      },
    ],
    eventName: 'updateMoveTaskOrderStatus',
    tableName: 'mto_service_items',
  };

  it('correctly matches the Create basic service item template', () => {
    const template = getTemplate(historyRecord);

    expect(template).toMatchObject(e);
    expect(template.getEventNameDisplay(template)).toEqual('Approved service item');
  });

  it('displays the expected details column value', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Move management')).toBeInTheDocument();
  });
});

import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/CancelMove/CancelMove';

describe('when given a Move cancellation history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'cancelMove',
    tableName: 'moves',
  };

  it('correctly matches the Move cancellation event to the proper template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the correct value in the details column', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Move canceled'));
  });
});

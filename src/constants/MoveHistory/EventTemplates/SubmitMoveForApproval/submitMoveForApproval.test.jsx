import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/SubmitMoveForApproval/submitMoveForApproval';

describe('When a customer signs and submits a move request', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: { status: 'NEEDS SERVICE COUNSELING' },
    eventName: 'submitMoveForApproval',
    tableName: 'moves',
  };

  it('matches to the proper template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays a confirmation message in the details column', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Received customer signature')).toBeInTheDocument();
  });

  it('displays event name', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay());
    expect(screen.getByText('Customer Signature')).toBeInTheDocument();
  });
});

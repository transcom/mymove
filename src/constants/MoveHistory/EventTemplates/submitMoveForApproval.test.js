import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/submitMoveForApproval';

describe('When a customer signs and submits a move request', () => {
  const item = {
    action: 'UPDATE',
    changedValues: { status: 'NEEDS SERVICE COUNSELING' },
    eventName: 'submitMoveForApproval',
    tableName: 'moves',
  };

  it('displays a confirmation message in the details column', () => {
    const result = getTemplate(item);

    expect(result).toMatchObject(e);
    render(result.getDetails(item));
    expect(screen.getByText('Received customer signature')).toBeInTheDocument();
  });
});

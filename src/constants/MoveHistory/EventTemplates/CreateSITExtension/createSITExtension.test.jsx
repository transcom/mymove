import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/CreateSITExtension/createSITExtension';

describe('when given a Create SIT extension item history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: { status: 'APPROVALS REQUESTED' },
    oldValues: { status: 'APPROVED' },
    eventName: 'createSITExtension',
    tableName: 'moves',
  };
  it('correctly matches to the Approved shipment, Updated move template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('correctly matches the Create SIT extension event', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': APPROVALS REQUESTED')).toBeInTheDocument();
  });
});

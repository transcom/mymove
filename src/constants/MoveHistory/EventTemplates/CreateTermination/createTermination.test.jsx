import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/CreateTermination/createTermination';

describe('when given a create termination history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: {
      status: 'TERMINATED_FOR_CAUSE',
      termination_comments: 'get in the choppuh',
      terminated_at: '2025-03-26 13:55:05.755',
    },
    eventName: 'createTermination',
    tableName: 'mto_shipments',
  };
  it('correctly matches to the template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('correctly matches the create termination event', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': TERMINATED FOR CAUSE')).toBeInTheDocument();
    expect(screen.getByText('Terminated at')).toBeInTheDocument();
    expect(screen.getByText(': 2025-03-26 13:55:05.755')).toBeInTheDocument();
    expect(screen.getByText('Comments')).toBeInTheDocument();
    expect(screen.getByText(': get in the choppuh')).toBeInTheDocument();
  });
});

import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import denySITExtensionMove from 'constants/MoveHistory/EventTemplates/DenySITExtension/denySITExtensionMove';
import Actions from 'constants/MoveHistory/Database/Actions';
import { MOVE_STATUSES } from 'shared/constants';

describe('when given a Deny SIT Extension Move move history record', () => {
  const historyRecord = {
    action: Actions.UPDATE,
    changedValues: {
      approved_at: '2025-04-09T13:43:45.206676+00:00',
      status: MOVE_STATUSES.APPROVED,
    },
    eventName: 'denySITExtension',
    tableName: 'moves',
  };

  it('matches the denySITExtensionMove template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(denySITExtensionMove);
  });

  it('returns the correct event display name', () => {
    expect(denySITExtensionMove.getEventNameDisplay()).toEqual('Updated move');
  });

  it('renders the move details via LabeledDetails', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));

    expect(screen.getByText('Approved at')).toBeInTheDocument();
    expect(screen.getByText(/2025-04-09T13:43:45.206676/)).toBeInTheDocument();
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(/APPROVED/)).toBeInTheDocument();
  });
});

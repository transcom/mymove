import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import approveSITExtensionMove from 'constants/MoveHistory/EventTemplates/ApproveSITExtension/approveSITExtensionMove';
import Actions from 'constants/MoveHistory/Database/Actions';
import { MOVE_STATUSES } from 'shared/constants';

describe('when given an Approve SIT Extension Move item history record', () => {
  const historyRecord = {
    action: Actions.UPDATE,
    changedValues: {
      approved_at: '2025-04-09T13:43:45.206676+00:00',
      status: MOVE_STATUSES.APPROVED,
    },
    eventName: 'approveSITExtension',
    tableName: 'moves',
  };

  it('correctly matches to the Approve SIT extension Move template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(approveSITExtensionMove);
  });

  it('returns the correct event display name', () => {
    expect(approveSITExtensionMove.getEventNameDisplay()).toEqual('Updated move');
  });

  it('renders the Approve SIT extension Move details correctly', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));

    expect(screen.getByText('Approved at')).toBeInTheDocument();
    expect(screen.getByText(/2025-04-09T13:43:45.206676/)).toBeInTheDocument();
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(/APPROVED/)).toBeInTheDocument();
  });
});

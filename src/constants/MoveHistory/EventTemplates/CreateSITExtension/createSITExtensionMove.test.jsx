import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import createSITExtensionMove from 'constants/MoveHistory/EventTemplates/CreateSITExtension/createSITExtensionMove';
import Actions from 'constants/MoveHistory/Database/Actions';
import { MOVE_STATUSES } from 'shared/constants';

describe('when given a Deny SIT Extension move history record', () => {
  const historyRecord = {
    action: Actions.UPDATE,
    changedValues: {
      approvals_requested_at: '2025-04-09T13:43:45.206676+00:00',
      status: MOVE_STATUSES.APPROVALS_REQUESTED,
    },
    eventName: 'createSITExtension',
    tableName: 'moves',
  };

  it('matches the createSITExtensionMove template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(createSITExtensionMove);
  });

  it('returns the correct event display name', () => {
    expect(createSITExtensionMove.getEventNameDisplay()).toEqual('Updated move');
  });

  it('renders the move details via LabeledDetails', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));

    expect(screen.getByText('Approvals requested at')).toBeInTheDocument();
    expect(screen.getByText(/2025-04-09T13:43:45.206676/)).toBeInTheDocument();
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(/APPROVALS REQUESTED/)).toBeInTheDocument();
  });
});

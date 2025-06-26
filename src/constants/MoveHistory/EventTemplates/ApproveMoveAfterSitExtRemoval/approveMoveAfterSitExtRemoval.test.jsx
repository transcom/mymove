import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import approveMoveAfterSitExtRemoval from 'constants/MoveHistory/EventTemplates/ApproveMoveAfterSitExtRemoval/approveMoveAfterSitExtRemoval';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import MOVE_STATUSES from 'constants/moves';

describe('when given a approveMoveAfterSitExtRemoval history record', () => {
  const historyRecord = {
    action: a.UPDATE,
    changedValues: {
      approved_at: '2025-06-06T19',
      status: MOVE_STATUSES.APPROVED,
    },
    context: [
      {
        eventName: 'updateMTOServiceItem',
        id: '16732a20-1e0e-470b-9f1f-7fd292a33059',
        objectId: '309bddbf-990c-4b92-9e18-698211d48c54"',
      },
    ],
    eventName: o.updateMTOServiceItem,
    tableName: t.moves,
  };

  it('matches the template from getTemplate', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(approveMoveAfterSitExtRemoval);
  });

  it('returns the correct event display name', () => {
    expect(approveMoveAfterSitExtRemoval.getEventNameDisplay()).toEqual('Updated move');
  });

  it('renders the details via LabeledDetails with merged changed values', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));

    // Check for the presence of the changed values.
    // The actual keys and values displayed depend on your LabeledDetails implementation.
    // Here we expect the values from changedValues to be rendered.
    expect(screen.getByText(/Status/i)).toBeInTheDocument();
    expect(screen.getByText(/Approved at/i)).toBeInTheDocument();
  });
});

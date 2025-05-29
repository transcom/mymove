import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/ReviewShipmentAddressUpdate/reviewShipmentAddressUpdateMove';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { MOVE_STATUSES } from 'shared/constants';

describe('when given a updated shipment address request, update move history record', () => {
  const historyRecord = {
    action: a.UPDATE,
    eventName: o.reviewShipmentAddressUpdate,
    tableName: t.moves,
    oldValues: { status: MOVE_STATUSES.APPROVALS_REQUESTED },
    changedValues: { status: MOVE_STATUSES.APPROVED, approved_at: '2025-04-09T13:43:45.206676+00:00' },
  };

  it('correctly matches the update service item status, update move event', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  describe('When given an updated shipment address request, update move history record', () => {
    it.each([
      ['Status', ': APPROVED'],
      ['Approved at', ': 2025-04-09T13:43:45.206676+00:00'],
    ])('displays the proper details value for %s', async (label, value) => {
      const template = getTemplate(historyRecord);
      render(template.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });

  it('displays correct details when a TOO is unassigned', () => {
    historyRecord.changedValues = {
      ...historyRecord.changedValues,
      too_assigned_id: null,
    };
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Task ordering officer unassigned')).toBeInTheDocument();
  });

  it('displays correct details when a TOO is unassigned and navigated from the destination request queue', () => {
    historyRecord.changedValues = {
      ...historyRecord.changedValues,
      too_destination_assigned_id: null,
    };
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Task ordering officer unassigned')).toBeInTheDocument();
  });
});

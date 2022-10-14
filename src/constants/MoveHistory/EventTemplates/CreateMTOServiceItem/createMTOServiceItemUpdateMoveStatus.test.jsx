import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import e from 'constants/MoveHistory/EventTemplates/CreateMTOServiceItem/createMTOServiceItemUpdateMoveStatus';

describe('when given a move status update with create mto service item history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: o.createMTOServiceItem,
    tableName: 'moves',
    detailsType: d.LABELED,
    oldValues: {
      status: 'Approved',
    },
  };
  const template = getTemplate(historyRecord);
  it('correctly matches the create MTO service item event', () => {
    expect(template).toMatchObject(e);
    expect(template.getEventNameDisplay()).toEqual('Updated move');
  });
  describe('when given a specific set of details', () => {
    it.each([['status', 'Approved']])('for label %s it displays the proper details value %s', async (label, value) => {
      render(template.getDetails(historyRecord));
      expect(screen.getByText(value, { exact: false })).toBeInTheDocument();
    });
  });
});

import { render, screen } from '@testing-library/react';

import a from 'constants/MoveHistory/Database/Actions';
import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import updateServiceMemberByCounselor from 'constants/MoveHistory/EventTemplates/UpdateServiceMember/updateServiceMemberByCounselor';

describe('When a service counselor updates shipping allowances', () => {
  const historyRecord = {
    action: a.UPDATE,
    tableName: t.service_members,
    eventName: o.counselingUpdateAllowance,
    eventNameDisplay: 'Updated profile',
    changedValues: {
      affiliation: 'NAVY',
      rank: 'E_4',
    },
  };
  it('it correctly matches the event that updates the service member profile ', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(updateServiceMemberByCounselor);
    expect(result.getEventNameDisplay()).toMatch(historyRecord.eventNameDisplay);
  });
  describe('it correctly renders the details component for the branch form', () => {
    it.each([
      ['Branch', ': Navy'],
      ['Rank', ': E-4'],
    ])('displays the correct details value for %s', async (label, value) => {
      const result = getTemplate(historyRecord);
      render(result.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import updateAllowanceServiceMemberByTOO from 'constants/MoveHistory/EventTemplates/UpdateServiceMember/updateServiceMemberByTOO';

describe('When a TOO updates shipping allowances', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: o.updateAllowance,
    tableName: t.service_members,
    eventNameDisplay: 'Updated profile',
    changedValues: {
      affiliation: 'AIR_FORCE',
      rank: 'E_2',
    },
  };
  it('correctly matches the update allowance event results in a change in service branch', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(updateAllowanceServiceMemberByTOO);
    expect(result.getEventNameDisplay()).toMatch(historyRecord.eventNameDisplay);
  });
  describe('it correctly displays the details component', () => {
    it.each([
      ['Branch', ': Air Force'],
      ['Rank', ': E-2'],
    ])('displays the correct details value for %s', async (label, value) => {
      const result = getTemplate(historyRecord);
      render(result.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import updateAllowanceServiceMemberByTOO from 'constants/MoveHistory/EventTemplates/UpdateServiceMember/updateServiceMemberByTOO';
import { ORDERS_PAY_GRADE_TYPE } from 'constants/orders';

describe('When a TOO updates shipping allowances', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updateAllowance',
    tableName: 'service_members',
    eventNameDisplay: 'Updated profile',
    changedValues: {
      affiliation: 'AIR_FORCE',
      grade: ORDERS_PAY_GRADE_TYPE.E_2,
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
      ['Pay grade', ': E-2'],
    ])('displays the correct details value for %s', async (label, value) => {
      const result = getTemplate(historyRecord);
      render(result.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

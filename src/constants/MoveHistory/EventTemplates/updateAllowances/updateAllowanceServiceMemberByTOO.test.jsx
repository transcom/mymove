import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import updateAllowanceServiceMemberByTOO from 'constants/MoveHistory/EventTemplates/updateAllowances/updateAllowanceServiceMemberByTOO';

describe('When a TOO updates shipping allowances', () => {
  const item = {
    action: 'UPDATE',
    eventName: o.updateAllowance,
    tableName: t.service_members,
    eventNameDisplay: 'Updated service member',
    changedValues: {
      affiliation: 'AIR_FORCE',
      rank: 'E_2',
    },
  };
  it('correctly matches the update allowance event results in a change in service branch', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(updateAllowanceServiceMemberByTOO);
    render(result.getDetails(item));
    expect(screen.getByText('Branch')).toBeInTheDocument();
    expect(screen.getByText(': Air Force')).toBeInTheDocument();
    expect(screen.getByText('Rank')).toBeInTheDocument();
    expect(screen.getByText(': E-2')).toBeInTheDocument();
  });
});

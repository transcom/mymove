import getTemplate from 'constants/MoveHistory/TemplateManager';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import updateAllowanceServiceMemberBranch from 'constants/MoveHistory/EventTemplates/updateAllowanceServiceMemberBranch';

describe('When a service counselor updates shipping allowances', () => {
  const item = {
    action: 'UPDATE',
    eventName: o.counselingUpdateAllowance,
    tableName: t.service_members,
    detailsType: d.LABELED,
    eventNameDisplay: 'Updated service member',
    changedValues: {
      affiliation: 'AIR_FORCE',
      rank: 'E_2',
    },
  };
  it('correctly matches the update allowance event results in a change in service branch', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(updateAllowanceServiceMemberBranch);
    expect(result.getEventNameDisplay()).toMatch(item.eventNameDisplay);
  });
});

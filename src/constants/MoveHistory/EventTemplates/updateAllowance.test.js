import getTemplate from 'constants/MoveHistory/TemplateManager';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import updateAllowance from 'constants/MoveHistory/EventTemplates/updateAllowance';

describe('When a service counselor updates shipping allowances', () => {
  const item = {
    action: 'UPDATE',
    eventName: o.updateAllowance,
    tableName: t.entitlements,
    detailsType: d.LABELED,
    eventNameDisplay: 'Updated allowances',
  };
  it('correctly matches the update allowances event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(updateAllowance);
  });
});

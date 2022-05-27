import o from 'constants/MoveHistory/UIDisplay/Operations';
import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/updateBillableWeightAsTIO';

describe('when given an update billable weights as tio history record', () => {
  const item = {
    action: 'UPDATE',
    changedValues: { authorized_weight: '7999' },
    eventName: o.updateBillableWeightAsTIO,
    tableName: 'entitlements',
  };
  it('correctly matches the update billable weights as tio event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay(item)).toEqual('Updated move');
  });
});

import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/updateOrder';

describe('when given an Order update history record', () => {
  const item = {
    action: 'UPDATE',
    eventName: 'updateOrder',
    tableName: 'orders',
    detailsType: d.LABELED,
    changedValues: { old_duty_location_id: 'ID1', new_duty_location_id: 'ID2' },
    context: [{ old_duty_location_name: 'old name', new_duty_location_name: 'new name' }],
  };
  it('correctly matches the Update orders event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    // expect to have merged context and changedValues
    expect(result.getDetailsLabeledDetails({ context: item.context, changedValues: item.changedValues })).toEqual({
      old_duty_location_id: 'ID1',
      new_duty_location_id: 'ID2',
      old_duty_location_name: 'old name',
      new_duty_location_name: 'new name',
    });
  });
});

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/updateReweighWeight';

describe('when given an updated reweigh weight', () => {
  const item = {
    action: 'UPDATE',
    changedValues: { weight: '9001' },
    context: [{ shipment_type: 'HHG' }],
    eventName: 'updateReweigh',
    tableName: 'reweighs',
  };
  it('correctly matches the update reweigh weight event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(
      result.getDetailsLabeledDetails({
        changedValues: item.changedValues,
        context: item.context,
      }),
    ).toEqual({
      shipment_type: 'HHG',
      reweigh_weight: '9001',
    });
  });
});

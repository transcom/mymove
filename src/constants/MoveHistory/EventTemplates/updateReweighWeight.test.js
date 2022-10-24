import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/updateReweighWeight';

describe('when given an updated reweigh weight', () => {
  const item = {
    action: 'UPDATE',
    changedValues: { weight: '9001' },
    context: [{ shipment_type: 'HHG', shipment_id_abbr: 'a1b2c' }],
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
      shipment_id_display: 'A1B2C',
      reweigh_weight: '9001',
    });
  });
});

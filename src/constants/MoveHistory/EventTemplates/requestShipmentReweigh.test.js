import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/requestShipmentReweigh';

describe('when given a Request shipment reweigh history record', () => {
  const item = {
    action: 'INSERT',
    context: [{ shipment_type: 'HHG', shipment_id_abbr: 'a1b2c' }],
    eventName: 'requestShipmentReweigh',
    tableName: 'reweighs',
  };
  it('correctly matches the Request shipment reweigh event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getDetailsPlainText(item)).toEqual('HHG shipment #A1B2C, reweigh requested');
  });
});

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateMTOShipment/createReweighRequestWeightUpdate';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import o from 'constants/MoveHistory/UIDisplay/Operations';

describe('when given a reweight request is created through weight update', () => {
  const item = {
    action: a.INSERT,
    eventName: o.updateMTOShipment,
    tableName: t.reweighs,
    context: [{ payment_request_number: '5650-7537-1', shipment_type: 'HHG', shipment_id_abbr: 'a1b2c' }],
  };
  it('correctly matches the create reweigh request weight update event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getDetailsPlainText(item)).toEqual('HHG shipment #A1B2C, reweigh requested');
  });
});

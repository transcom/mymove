import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/updateServiceItemStatus';

describe('when given an Approved service item history record', () => {
  const item = {
    action: 'UPDATE',
    changedValues: { status: 'APPROVED' },
    context: [{ name: 'Domestic origin price', shipment_type: 'HHG_INTO_NTS_DOMESTIC' }],
    eventName: 'updateMTOServiceItemStatus',
    tableName: 'mto_service_items',
  };
  it('correctly matches the Approved service item event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getDetailsPlainText(item)).toEqual('NTS shipment, Domestic origin price');
  });
});

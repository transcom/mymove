import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/deleteShipment';

describe('When a service counselor deletes a shipment', () => {
  const item = {
    action: 'UPDATE',
    eventName: 'deleteShipment',
    oldValues: { shipment_type: 'HHG' },
    tableName: 'mto_shipments',
  };
  it('correctly displays a message indicating a HHG shipment was deleted', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getDetailsPlainText(item)).toEqual('HHG shipment deleted');
  });
});

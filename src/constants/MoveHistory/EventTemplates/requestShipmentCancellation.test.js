import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/requestShipmentCancellation';

describe('when given a Request shipment cancellation history record', () => {
  const item = {
    action: 'UPDATE',
    changedValues: {
      status: 'DRAFT',
    },
    eventName: 'requestShipmentCancellation',
    oldValues: { shipment_type: 'PPM' },
    tableName: '',
  };
  it('correctly matches the Request shipment cancellation event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getDetailsPlainText(item)).toEqual('Requested cancellation for PPM shipment');
  });
});

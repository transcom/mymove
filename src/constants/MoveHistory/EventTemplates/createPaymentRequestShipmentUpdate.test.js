import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/createPaymentRequestShipmentUpdate';

describe('when given a payment request is created through shipment update', () => {
  const item = {
    action: 'INSERT',
    eventName: 'updateMTOShipment',
    tableName: 'payment_requests',
  };
  it('correctly matches the Request shipment reweigh event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getStatusDetails(item)).toEqual('Pending');
  });
});

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/createPaymentRequestReweighUpdate';

describe('when given a payment request is created through reweigh', () => {
  const item = {
    action: 'INSERT',
    eventName: 'updateReweigh',
    tableName: 'payment_requests',
  };
  it('correctly matches the Request shipment reweigh event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getStatusDetails(item)).toEqual('Pending');
  });
});

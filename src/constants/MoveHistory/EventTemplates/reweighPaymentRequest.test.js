import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import e from 'constants/MoveHistory/EventTemplates/reweighPaymentRequest';

describe('reweighs update', () => {
  const item = {
    action: 'UPDATE',
    eventName: o.updateReweigh,
    tableName: 'reweighs',
    context: [{ shipment_type: 'HHG', payment_request_number: '0288-7994-1' }],
  };
  it('correctly matches the reweigh payment request', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getStatusDetails(item)).toEqual('Recalculated payment request');
  });
});

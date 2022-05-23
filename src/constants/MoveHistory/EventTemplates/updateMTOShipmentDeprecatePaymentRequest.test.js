import o from 'constants/MoveHistory/UIDisplay/Operations';
import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/updateMTOShipmentDeprecatePaymentRequest';

describe('when given a deprecated payment request history record', () => {
  const item = {
    action: 'UPDATE',
    eventName: o.updateMTOShipment,
    tableName: 'payment_requests',
    changedValues: {
      status: 'DEPRECATED',
    },
  };
  it('correctly matches the deprecated payment request', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
  });
});

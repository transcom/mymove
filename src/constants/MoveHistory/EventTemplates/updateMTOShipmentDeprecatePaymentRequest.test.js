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

describe('when given a recalculation_of_payment_request_id, displays recalculated payment request', () => {
  const item = {
    action: 'UPDATE',
    eventName: o.updateMTOShipment,
    tableName: 'payment_requests',
    changedValues: {
      recalculation_of_payment_request_id: '1234-5789-1',
    },
  };
  it('correctly matches the recalculated payment request', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getStatusDetails(item)).toEqual('Recalculated payment request');
  });
});

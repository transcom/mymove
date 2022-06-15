import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/updatePaymentRequest';

describe('when a payment request has an update', () => {
  const item = {
    action: 'UPDATE',
    tableName: 'payment_requests',
  };
  it('correctly matches the update payment request event for when a payment has been sent to GEX', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(
      result.getStatusDetails({
        changedValues: { status: 'SENT_TO_GEX' },
      }),
    ).toEqual('Sent to GEX');
  });
  it('correctly matches the update payment request event for when a payment has been received by GEX', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(
      result.getStatusDetails({
        changedValues: { status: 'RECEIVED_BY_GEX' },
      }),
    ).toEqual('Received');
  });

  it('correctly matches the update payment request event for when theres and EDI error', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(
      result.getStatusDetails({
        changedValues: { status: 'EDI_ERROR' },
      }),
    ).toEqual('EDI error');
  });
});

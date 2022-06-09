import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/createPaymentRequest';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import o from 'constants/MoveHistory/UIDisplay/Operations';

describe('when given a payment request is created through shipment update', () => {
  const item = {
    action: a.INSERT,
    eventName: o.createPaymentRequest,
    tableName: t.payment_requests,
    context: [
      {
        name: 'Test service',
        price: '10123',
        status: 'REQUESTED',
        shipment_id: '123',
        shipment_type: 'HHG',
      },
      {
        name: 'Domestic uncrating',
        price: '5555',
        status: 'REQUESTED',
        shipment_id: '456',
        shipment_type: 'HHG_INTO_NTS_DOMESTIC',
      },
      { name: 'Move management', price: '1234', status: 'REQUESTED' },
    ],
  };
  it('correctly matches the create payment request event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getLabeledPaymentRequestDetails(item.context)).toMatchObject({
      moveServices: 'Move management',
      shipmentServices: [
        { serviceItems: 'Test service', shipmentId: '123', shipmentType: 'HHG' },
        { serviceItems: 'Domestic uncrating', shipmentId: '456', shipmentType: 'HHG_INTO_NTS_DOMESTIC' },
      ],
    });
  });
});

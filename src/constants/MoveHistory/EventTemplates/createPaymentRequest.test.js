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
      { name: 'Test Service', price: '', status: 'REQUESTED' },
      { name: 'Domestic origin price', price: '13550', status: 'REQUESTED' },
    ],
  };
  it('correctly matches the create payment request event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getDetailsLabeledDetails(item)).toMatchObject({
      move_services: '',
      shipment_services: 'Test Service, Domestic origin price',
      shipment_type: 'HHG',
    });
  });
});

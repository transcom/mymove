import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/approveShipmentDiversion';

describe('when given an Approved shipment diversion history record', () => {
  const item = {
    changedValues: { status: 'APPROVED' },
    eventName: 'approveShipmentDiversion',
    oldValues: { shipment_type: 'HHG' },
  };
  it('correctly matches the Approved shipment event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getDetailsPlainText(item)).toEqual('HHG shipment');
  });
});

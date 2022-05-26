import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/createStandardServiceItem';

describe('when given a Create standard service item history record', () => {
  const item = {
    action: 'INSERT',
    context: [
      {
        shipment_type: 'HHG',
        name: 'Domestic linehaul',
      },
    ],
    eventName: 'approveShipment',
    tableName: 'mto_service_items',
  };
  it('correctly matches the Create standard service item event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay(result)).toEqual('Approved service item');
    expect(result.getDetailsPlainText(item)).toEqual('HHG shipment, Domestic linehaul');
  });
});

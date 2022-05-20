import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/createMTOShipment';

describe('when given a Create MTO Shipment event', () => {
  const item = {
    action: 'INSERT',
    context: [
      {
        name: 'Move management',
      },
    ],
    eventName: 'createMTOShipment',
    tableName: 'mto_shipments',
  };
  it('correctly matches the Create MTO Shipment event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay(result)).toEqual('Submitted/Requested shipments');
    expect(result.getDetailsPlainText(item)).toEqual('-');
  });
});

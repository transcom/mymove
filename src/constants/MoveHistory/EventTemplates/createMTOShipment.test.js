import getTemplate from 'constants/MoveHistory/TemplateManager';
import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import e from 'constants/MoveHistory/EventTemplates/createMTOShipment';

describe('when given a Create mto shipment history record', () => {
  const item = {
    action: a.INSERT,
    eventName: o.createMTOShipment,
    tableName: t.mto_shipments,
  };
  it('correctly matches the Create basic service item event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay(result)).toEqual('Created shipment');
  });
});

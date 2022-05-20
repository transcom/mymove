import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/undefined';

describe('when given an unidentifiable move history record', () => {
  const item = {
    action: 'UPDATE',
    eventName: 'testEventName',
    tableName: 'mto_agents',
  };
  it('correctly matches the Undefined move history event', () => {
    const result = getTemplate(item);
    expect(result).toEqual(e);
    expect(result.getEventNameDisplay(item)).toEqual('Updated shipment');
  });
});

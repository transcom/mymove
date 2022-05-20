import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/createBasicServiceItem';

describe('when given a Create basic service item history record', () => {
  const item = {
    action: 'INSERT',
    context: [
      {
        name: 'Move management',
      },
    ],
    eventName: 'updateMoveTaskOrderStatus',
    tableName: 'mto_service_items',
  };
  it('correctly matches the Create basic service item event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay(result)).toEqual('Approved service item');
    expect(result.getDetailsPlainText(item)).toEqual('Move management');
  });
});

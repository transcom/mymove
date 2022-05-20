import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/createOrders';

describe('when given a Submitted orders history record', () => {
  const item = {
    action: 'INSERT',
    changedValues: {
      status: 'DRAFT',
    },
    eventName: 'createOrders',
    tableName: 'orders',
  };
  it('correctly matches the Submitted orders event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
  });
});

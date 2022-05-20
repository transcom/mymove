import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/updateMoveTaskOrderStatus';

describe('when given a Move approved history record', () => {
  const item = {
    action: 'UPDATE',
    changedValues: {
      available_to_prime_at: '2022-04-13T15:21:31.746028+00:00',
      status: 'APPROVED',
    },
    eventName: 'updateMoveTaskOrderStatus',
    tableName: 'moves',
  };
  it('correctly matches the Update move task order status event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay(item)).toEqual('Approved move');
    expect(result.getDetailsPlainText(item)).toEqual('Created Move Task Order (MTO)');
  });
});

describe('when given a Move status update history record', () => {
  const item = {
    action: 'UPDATE',
    changedValues: { status: 'CANCELED' },
    eventName: 'updateMoveTaskOrderStatus',
    tableName: 'moves',
  };
  it('correctly matches the Update move task order status event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay(item)).toEqual('Move status updated');
    expect(result.getDetailsPlainText(item)).toEqual('-');
  });
});

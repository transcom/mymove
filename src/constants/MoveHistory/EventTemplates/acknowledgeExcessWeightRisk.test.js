import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/acknowledgeExcessWeightRisk';

describe('when given an Acknowledge excess weight risk history record', () => {
  const item = {
    action: 'UPDATE',
    eventName: 'acknowledgeExcessWeightRisk',
    tableName: 'moves',
  };
  it('correctly matches the Acknowledge excess weight risk event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getDetailsPlainText(item)).toEqual('Dismissed excess weight alert');
  });
});

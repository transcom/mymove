import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/setFinancialReviewFlag';

describe('when given a Set financial review flag event for flagged move history record', () => {
  const item = {
    action: 'UPDATE',
    eventName: 'setFinancialReviewFlag',
    changedValues: { financial_review_flag: 'true' },
    tableName: 'moves',
  };
  it('correctly matches the Set financial review flag event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getDetailsPlainText(item)).toEqual('Move flagged for financial review');
  });
});

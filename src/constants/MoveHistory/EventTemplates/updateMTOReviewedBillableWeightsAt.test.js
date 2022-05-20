import updateMTOReviewedBillableWeightsAt from './updateMTOReviewedBillableWeightsAt';

import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('when given an MTO Reviewed Billable Weight At event', () => {
  const item = {
    action: 'UPDATE',
    eventName: 'UpdateMTOReviewedBillableWeightsAt',
    tableName: 'moves',
  };
  it('correctly matches the MTO Reviewed Billable Weight At event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(updateMTOReviewedBillableWeightsAt);
    expect(result.getDetailsPlainText(item)).toEqual('Reviewed weights');
  });
});

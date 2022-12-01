import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/SetFinancialReviewFlag/setFinancialReviewFlag';

describe('when given a Set financial review flag event for flagged move history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'setFinancialReviewFlag',
    changedValues: {
      financial_review_flag: 'true',
      financial_review_remarks: 'This shipment needs to be reviewed asap',
    },
    tableName: 'moves',
  };

  it('correctly matches the Set financial review flag event', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('correctly displays changed values in details column', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Move flagged for financial review')).toBeInTheDocument();
  });

  describe('When given a specific set of details', () => {
    it.each([['Financial review remarks', ': This shipment needs to be reviewed asap']])(
      'displays the proper details value for %s',
      async (label, value) => {
        const template = getTemplate(historyRecord);

        render(template.getDetails(historyRecord));
        expect(screen.getByText(label)).toBeInTheDocument();
        expect(screen.getByText(value)).toBeInTheDocument();
      },
    );
  });
});

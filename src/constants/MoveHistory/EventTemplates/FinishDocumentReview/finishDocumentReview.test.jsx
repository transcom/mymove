import { screen, render } from '@testing-library/react';

import e from 'constants/MoveHistory/EventTemplates/FinishDocumentReview/finishDocumentReview';
import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('When given a completed services counseling for a move', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'finishDocumentReview',
    tableName: 'ppm_shipments',
    context: [
      {
        shipment_id_abbr: 'acf7b',
        shipment_locator: 'RQ38D4-01',
        shipment_type: 'PPM',
      },
    ],
    changedValues: { status: 'CLOSEOUT_COMPLETE' },
  };
  it('correctly matches the update mto status services counseling completed event to the proper template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper header for the given history record', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('PPM shipment #RQ38D4-01')).toBeInTheDocument();
  });

  describe('When given a specific set of details for a cancelled shipment', () => {
    it.each([['Status', ': PAYMENT APPROVED']])('displays the proper details value for %s', async (label, value) => {
      const template = getTemplate(historyRecord);
      render(template.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

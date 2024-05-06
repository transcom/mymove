import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { ppmExpenseTypes } from 'constants/ppmExpenseTypes';

const expenseTypes = ppmExpenseTypes.map((expense) => [expense.value, expense.key]);

describe('When given a created moving expense history record', () => {
  const historyRecord = {
    action: a.INSERT,
    changedValues: {
      amount: null,
      deleted_at: null,
      description: null,
      document_id: 'fde71ebb-fbf0-450d-92f1-c0f9344f9c84',
      id: 'bbe57f42-9a37-42a2-9594-503bb5c09022',
      missing_receipt: null,
      moving_expense_type: null,
      paid_with_gtcc: null,
      ppm_shipment_id: '87757854-3c57-4aaf-a2f3-0ae701c2bb0a',
      reason: null,
      sit_end_date: null,
      sit_start_date: null,
      status: null,
    },
    oldValues: {},
    context: [
      {
        shipment_id_abbr: '125d1',
        shipment_locator: 'RQ38D4-01',
        shipment_type: 'PPM',
      },
    ],
    eventName: o.createMovingExpense,
    tableName: t.moving_expenses,
  };

  it('displays event properly', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    expect(screen.getByText('Created moving expense')).toBeInTheDocument();
  });

  describe('properly renders shipment labels for ', () => {
    it.each(expenseTypes)('%s receipts', (label, docType) => {
      historyRecord.oldValues.moving_expense_type = docType;
      const template = getTemplate(historyRecord);

      render(template.getDetails(historyRecord));
      expect(screen.getByText(`PPM shipment #RQ38D4-01, ${label}`)).toBeInTheDocument();
    });
  });
});

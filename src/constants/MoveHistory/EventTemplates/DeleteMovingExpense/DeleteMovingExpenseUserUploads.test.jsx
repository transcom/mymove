import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { ppmExpenseTypes } from 'constants/ppmExpenseTypes';

const expenseTypes = ppmExpenseTypes.map((expense) => [expense.value, expense.key]);

describe('When given a deleted expense receipt upload', () => {
  const historyRecord = {
    action: a.UPDATE,
    changedValues: {
      deleted_at: '2024-02-15T08:41:06.592578+00:00',
    },
    oldValues: {},
    context: [
      {
        filename: 'expense.png',
        shipment_id_abbr: '7f559',
        shipment_locator: 'RQ38D4-01',
        shipment_type: 'PPM',
        upload_type: 'expenseReceipt',
      },
    ],
    eventName: o.deleteMovingExpense,
    tableName: t.user_uploads,
  };

  it('displays event properly', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    expect(screen.getByText('Deleted document')).toBeInTheDocument();
  });

  describe('properly renders shipment labels for ', () => {
    it.each(expenseTypes)('%s receipts', (label, docType) => {
      historyRecord.oldValues.moving_expense_type = docType;
      const template = getTemplate(historyRecord);

      render(template.getDetails(historyRecord));
      expect(screen.getByText(`PPM shipment #RQ38D4-01, ${label}`)).toBeInTheDocument();
    });
  });

  it('displays details of a deleted expense document', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Document type')).toBeInTheDocument();
    expect(screen.getByText(': Expense receipt')).toBeInTheDocument();
    expect(screen.getByText('Filename')).toBeInTheDocument();
    expect(screen.getByText(': expense.png')).toBeInTheDocument();
  });
});

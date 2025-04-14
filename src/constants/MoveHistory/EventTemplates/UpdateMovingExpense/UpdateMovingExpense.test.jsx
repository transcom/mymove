import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { ppmExpenseTypes } from 'constants/ppmExpenseTypes';
import { PPM_DOCUMENT_STATUS } from 'constants/ppms';

const expenseTypes = ppmExpenseTypes.map((expense) => [expense.value, expense.key]);

describe('When given an updated expense document it', () => {
  const expenseRecord = {
    action: a.UPDATE,
    changedValues: {},
    oldValues: {},
    context: [
      {
        shipment_id_abbr: '71f6f',
        shipment_locator: 'RQ38D4-01',
        shipment_type: 'PPM',
      },
    ],
    eventName: o.updateMovingExpense,
    tableName: t.moving_expenses,
  };

  describe('properly renders approvals and amount changes to ', () => {
    it.each(expenseTypes)('%s documents', (label, docType) => {
      expenseRecord.oldValues.moving_expense_type = docType;
      expenseRecord.changedValues = {
        status: PPM_DOCUMENT_STATUS.APPROVED,
        amount: '999999',
      };

      const template = getTemplate(expenseRecord);

      render(template.getDetails(expenseRecord));
      expect(screen.getByText(`PPM shipment #RQ38D4-01, ${label}`)).toBeInTheDocument();
      expect(screen.getByText(': APPROVED')).toBeInTheDocument();
      expect(screen.getByText(': $9,999.99')).toBeInTheDocument();
    });
  });

  describe('properly renders rejections and remarks for ', () => {
    it.each(expenseTypes)('%s documents', (label, docType) => {
      expenseRecord.oldValues.moving_expense_type = docType;
      expenseRecord.changedValues = {
        status: PPM_DOCUMENT_STATUS.REJECTED,
        reason: 'cannot read document',
      };

      const template = getTemplate(expenseRecord);

      render(template.getDetails(expenseRecord));
      expect(screen.getByText(`PPM shipment #RQ38D4-01, ${label}`)).toBeInTheDocument();
      expect(screen.getByText(': REJECTED')).toBeInTheDocument();
      expect(screen.getByText(': cannot read document')).toBeInTheDocument();
    });
  });

  describe('properly renders excluded expenses and remarks for ', () => {
    it.each(expenseTypes)('%s documents', (label, docType) => {
      expenseRecord.oldValues.moving_expense_type = docType;
      expenseRecord.changedValues = {
        status: PPM_DOCUMENT_STATUS.EXCLUDED,
        reason: 'claim on taxes',
      };

      const template = getTemplate(expenseRecord);

      render(template.getDetails(expenseRecord));
      expect(screen.getByText(`PPM shipment #RQ38D4-01, ${label}`)).toBeInTheDocument();
      expect(screen.getByText(': EXCLUDED')).toBeInTheDocument();
      expect(screen.getByText(': claim on taxes')).toBeInTheDocument();
    });
  });

  describe('properly renders updated expenses', () => {
    it.each(expenseTypes)('%s documents', (label, docType) => {
      expenseRecord.oldValues.moving_expense_type = docType;
      expenseRecord.changedValues = {
        missing_receipt: true,
        moving_expense_type: docType,
        paid_with_gtcc: false,
      };
      const template = getTemplate(expenseRecord);

      render(template.getDetails(expenseRecord));
      expect(screen.getByText(`PPM shipment #RQ38D4-01, ${label}`)).toBeInTheDocument();
      expect(screen.getByText('Paid with gtcc')).toBeInTheDocument();
      expect(screen.getByText(': No')).toBeInTheDocument();
      expect(screen.getByText('Missing receipt')).toBeInTheDocument();
      expect(screen.getByText(': Yes')).toBeInTheDocument();
      expect(screen.getByText('Expense type')).toBeInTheDocument();
      expect(
        screen.getByText((content) =>
          content.replace(/\s+/g, ' ').trim().toLowerCase().includes(`: ${label}`.toLowerCase()),
        ),
      ).toBeInTheDocument();
    });
  });

  it('handles SIT dates', () => {
    expenseRecord.changedValues = {
      sit_end_date: '2024-01-01',
      sit_start_date: '2023-12-12',
    };

    const template = getTemplate(expenseRecord);
    render(template.getDetails(expenseRecord));
    expect(screen.getByText('SIT start date')).toBeInTheDocument();
    expect(screen.getByText(': 12 Dec 2023')).toBeInTheDocument();
    expect(screen.getByText('SIT end date')).toBeInTheDocument();
    expect(screen.getByText(': 01 Jan 2024')).toBeInTheDocument();
  });
});

import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { ppmExpenseTypes } from 'constants/ppmExpenseTypes';
import ppms from 'constants/ppms';

const expenseTypes = ppmExpenseTypes.map((expense) => [expense.value, expense.key]);

describe('When given an updated expense document it', () => {
  const expenseRecord = {
    action: a.UPDATE,
    changedValues: {},
    context: [
      {
        shipment_id_abbr: '71f6f',
        shipment_type: 'PPM',
      },
    ],
    eventName: o.updateMovingExpense,
    tableName: t.moving_expenses,
  };

  describe('properly renders approvals and amount changes to ', () => {
    it.each(expenseTypes)('%s documents', (label, docType) => {
      expenseRecord.context[0].moving_expense_type = docType;
      expenseRecord.changedValues = {
        status: ppms.APPROVED,
        amount: '999999',
      };

      const template = getTemplate(expenseRecord);

      render(template.getDetails(expenseRecord));
      expect(screen.getByText(`PPM shipment #71F6F, ${label}`)).toBeInTheDocument();
      expect(screen.getByText(': APPROVED')).toBeInTheDocument();
      expect(screen.getByText(': $9,999.99')).toBeInTheDocument();
    });
  });

  describe('properly renders rejections and remarks for ', () => {
    it.each(expenseTypes)('%s documents', (label, docType) => {
      expenseRecord.context[0].moving_expense_type = docType;
      expenseRecord.changedValues = {
        status: ppms.REJECTED,
        reason: 'cannot read document',
      };

      const template = getTemplate(expenseRecord);

      render(template.getDetails(expenseRecord));
      expect(screen.getByText(`PPM shipment #71F6F, ${label}`)).toBeInTheDocument();
      expect(screen.getByText(': REJECTED')).toBeInTheDocument();
      expect(screen.getByText(': cannot read document')).toBeInTheDocument();
    });
  });

  describe('properly renders excluded expenses and remarks for ', () => {
    it.each(expenseTypes)('%s documents', (label, docType) => {
      expenseRecord.context[0].moving_expense_type = docType;
      expenseRecord.changedValues = {
        status: ppms.EXCLUDED,
        reason: 'claim on taxes',
      };

      const template = getTemplate(expenseRecord);

      render(template.getDetails(expenseRecord));
      expect(screen.getByText(`PPM shipment #71F6F, ${label}`)).toBeInTheDocument();
      expect(screen.getByText(': EXCLUDED')).toBeInTheDocument();
      expect(screen.getByText(': claim on taxes')).toBeInTheDocument();
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

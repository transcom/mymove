import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { ppmExpenseTypes } from 'constants/ppmExpenseTypes';
import { PPM_UPLOAD_TYPES, PPM_UPLOAD_TYPES_LABELS } from 'constants/ppmUploadTypes';

const expenseTypes = ppmExpenseTypes.map((expense) => [expense.value, expense.key]);
const uploadTypes = PPM_UPLOAD_TYPES.map((uploadType) => [PPM_UPLOAD_TYPES_LABELS[uploadType], uploadType]);

describe('When given a created pro-gear set history record', () => {
  const historyRecord = {
    action: a.INSERT,
    changedValues: {
      belongs_to_self: null,
      deleted_at: null,
      description: null,
      document_id: '4a16b7ef-7e1b-44dc-99e7-75cbf3141da2',
      has_weight_tickets: null,
      id: '6cf4c8ff-c0ab-4d7c-8330-3ea03e7a3f9a',
      ppm_shipment_id: '87757854-3c57-4aaf-a2f3-0ae701c2bb0a',
      reason: null,
      status: null,
      weight: null,
    },
    oldValues: {},
    context: [
      {
        filename: 'filename.png',
        shipment_id_abbr: '125d1',
        shipment_locator: 'RQ38D4-01',
        shipment_type: 'PPM',
      },
    ],
    eventName: o.createPPMUpload,
    tableName: t.user_uploads,
  };

  it('displays event properly', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    expect(screen.getByText('Uploaded document')).toBeInTheDocument();
  });

  describe('properly renders the details of ', () => {
    it.each(uploadTypes)('%s', (label, uploadType) => {
      historyRecord.context[0].upload_type = uploadType;
      const template = getTemplate(historyRecord);

      render(template.getDetails(historyRecord));
      expect(screen.getByText('Document type')).toBeInTheDocument();
      expect(screen.getByText(`: ${label}`)).toBeInTheDocument();
      expect(screen.getByText('Filename')).toBeInTheDocument();
      expect(screen.getByText(`: filename.png`)).toBeInTheDocument();
    });
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

import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

describe('When given a deleted trip weight ticket upload', () => {
  const historyRecord = {
    action: a.UPDATE,
    changedValues: {
      deleted_at: '2024-02-15T08:41:06.592578+00:00',
    },
    context: [
      {
        filename: 'pgWeight.jpg',
        moving_expense_type: '',
        shipment_id_abbr: '125d1',
        shipment_locator: 'RQ38D4-01',
        shipment_type: 'PPM',
      },
    ],
    eventName: o.deleteWeightTicket,
    tableName: t.user_uploads,
  };

  it('displays event properly', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    expect(screen.getByText('Deleted document')).toBeInTheDocument();
  });

  it('displays trip event properly', () => {
    historyRecord.context[0].upload_type = 'fullWeightTicket';
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    expect(screen.getByText('Deleted trip document')).toBeInTheDocument();
  });

  it('displays details of shipment type, shipment ID', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('PPM shipment #RQ38D4-01')).toBeInTheDocument();
  });

  describe('displays details of a deleted ', () => {
    it('empty weight ticket', () => {
      historyRecord.context[0].upload_type = 'emptyWeightTicket';
      const template = getTemplate(historyRecord);

      render(template.getDetails(historyRecord));
      expect(screen.getByText('Document type')).toBeInTheDocument();
      expect(screen.getByText(': Empty weight ticket')).toBeInTheDocument();
      expect(screen.getByText('Filename')).toBeInTheDocument();
      expect(screen.getByText(': pgWeight.jpg')).toBeInTheDocument();
    });
    it('full weight ticket', () => {
      historyRecord.context[0].upload_type = 'fullWeightTicket';
      const template = getTemplate(historyRecord);

      render(template.getDetails(historyRecord));
      expect(screen.getByText('Document type')).toBeInTheDocument();
      expect(screen.getByText(': Full weight ticket')).toBeInTheDocument();
      expect(screen.getByText('Filename')).toBeInTheDocument();
      expect(screen.getByText(': pgWeight.jpg')).toBeInTheDocument();
    });
    it('trailer weight ticket', () => {
      historyRecord.context[0].upload_type = 'trailerWeightTicket';
      const template = getTemplate(historyRecord);

      render(template.getDetails(historyRecord));
      expect(screen.getByText('Document type')).toBeInTheDocument();
      expect(screen.getByText(': Trailer weight ticket')).toBeInTheDocument();
      expect(screen.getByText('Filename')).toBeInTheDocument();
      expect(screen.getByText(': pgWeight.jpg')).toBeInTheDocument();
    });
  });
});

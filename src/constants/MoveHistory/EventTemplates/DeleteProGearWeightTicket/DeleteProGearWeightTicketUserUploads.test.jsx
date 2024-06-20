import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

describe('When given a deleted pro-gear weight ticket upload', () => {
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
        upload_type: 'spouseProGearWeightTicket',
      },
    ],
    eventName: o.deleteProGearWeightTicket,
    tableName: t.user_uploads,
  };

  it('displays event properly', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    expect(screen.getByText('Deleted document')).toBeInTheDocument();
  });

  it('displays details of shipment type, shipment ID', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('PPM shipment #RQ38D4-01, Spouse pro-gear')).toBeInTheDocument();
  });

  it('displays details the deleted document', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Document type')).toBeInTheDocument();
    expect(screen.getByText(': Spouse pro-gear weight ticket')).toBeInTheDocument();
    expect(screen.getByText('Filename')).toBeInTheDocument();
    expect(screen.getByText(': pgWeight.jpg')).toBeInTheDocument();
  });
});

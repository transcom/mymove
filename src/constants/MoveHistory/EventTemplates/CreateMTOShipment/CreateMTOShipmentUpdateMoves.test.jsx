import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

describe('When given a move with an updated ppm type', () => {
  const historyRecord = {
    action: a.UPDATE,
    changedValues: {
      ppm_type: 'FULL',
      closeout_office_name: 'PPPO Scott AFB - USAF',
    },
    context: [
      {
        shipment_id_abbr: '125d1',
        shipment_locator: 'RQ38D4-01',
        shipment_type: 'PPM',
      },
    ],
    eventName: o.createMTOShipment,
    tableName: t.moves,
  };

  it('displays event properly', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    expect(screen.getByText('Updated move')).toBeInTheDocument();
  });

  it('displays details ppmType', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('PPM type')).toBeInTheDocument();
    expect(screen.getByText(': FULL')).toBeInTheDocument();
  });
});

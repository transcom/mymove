import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateMTOShipment/patchMTOShipment';

describe('when given a patch move record', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'patchMove',
    tableName: 'moves',
    changedValues: { closeout_office_name: 'PPPO Scott AFB - USAF' },
  };

  it('correctly matches the patch move event', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the proper update order record', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Closeout office')).toBeInTheDocument();
    expect(screen.getByText(': PPPO Scott AFB - USAF')).toBeInTheDocument();
  });
});

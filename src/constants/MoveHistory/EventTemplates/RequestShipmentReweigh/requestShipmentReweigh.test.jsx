import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/RequestShipmentReweigh/requestShipmentReweigh';

describe('when given a Request shipment reweigh history record', () => {
  const historyRecord = {
    action: 'INSERT',
    context: [{ shipment_type: 'HHG', shipment_id_abbr: 'a1b2c', shipment_locator: 'ABC123-01' }],
    eventName: 'requestShipmentReweigh',
    tableName: 'reweighs',
  };

  it('correctly matches the Request shipment reweigh to the proper template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('renders the reweigh details message on the screen', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('HHG shipment #ABC123-01, reweigh requested')).toBeInTheDocument();
  });
});

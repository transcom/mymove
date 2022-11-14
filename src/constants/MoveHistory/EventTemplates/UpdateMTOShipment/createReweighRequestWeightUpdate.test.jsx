import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import createReweighRequestWeightUpdate from 'constants/MoveHistory/EventTemplates/UpdateMTOShipment/createReweighRequestWeightUpdate';

describe('when given an mto shipment reweigh request', () => {
  const historyRecord = {
    action: 'INSERT',
    eventName: 'updateMTOShipment',
    tableName: 'reweighs',
    context: [
      {
        shipment_type: 'PPM',
        shipment_id_abbr: 'b4b4b',
      },
    ],
  };
  it('correctly matches the createReweighRequestWeightUpdate event', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(createReweighRequestWeightUpdate);
    expect(result.getEventNameDisplay()).toMatch('Updated shipment');
  });
  it('displays the correct details component', () => {
    const result = getTemplate(historyRecord);
    render(result.getDetails(historyRecord));
    expect(screen.getByText('PPM shipment #B4B4B, reweigh requested')).toBeInTheDocument();
  });
});

import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import updateMTOShipmentPrimeEstimatedExcessWeight from 'constants/MoveHistory/EventTemplates/UpdateMTOShipment/updateMTOShipmentPrimeEstimatedExcessWeight';

describe("when Prime user updates a shipment's estimated to one that exceeds the weight limit", () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updateMTOShipment',
    tableName: 'moves',
    eventNameDisplay: 'Updated move',
    changedValues: {
      excess_weight_qualified_at: '2022-09-07T21:26:47.833768+00:00',
    },
  };
  it('correctly matches the Update mto shipment event', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(updateMTOShipmentPrimeEstimatedExcessWeight);
    expect(result.getEventNameDisplay()).toEqual('Updated move');
  });
  describe('it displays the proper labeled details for the component', () => {
    const result = getTemplate(historyRecord);
    render(result.getDetails(historyRecord));
    expect(
      screen.getByText('Flagged for excess weight, total estimated weight > 90% weight allowance'),
    ).toBeInTheDocument();
  });
});

import getTemplate from 'constants/MoveHistory/TemplateManager';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import updateMTOShipmentPrimeEstimatedExcessWeight from 'constants/MoveHistory/EventTemplates/updateMTOShipmentPrimeEstimatedExcessWeight';

describe("when Prime user updates a shipment's estimated to one that exceeds the weight limit", () => {
  const item = {
    action: 'UPDATE',
    eventName: o.updateMTOShipment,
    tableName: 'moves',
    detailsType: d.PLAIN_TEXT,
    eventNameDisplay: 'Updated move',
    changedValues: {
      excess_weight_qualified_at: '2022-09-07T21:26:47.833768+00:00',
    },
  };
  it('correctly matches the Update mto shipment event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(updateMTOShipmentPrimeEstimatedExcessWeight);
    expect(result.getDetailsPlainText(item)).toEqual(
      'Flagged for excess weight, total estimated weight > 90% weight allowance',
    );
  });
});

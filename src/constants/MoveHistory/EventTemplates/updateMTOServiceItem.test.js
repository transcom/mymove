import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import e from 'constants/MoveHistory/EventTemplates/updateMTOServiceItem';

describe('when given a Update basic service item history record', () => {
  const item = {
    action: a.UPDATE,
    changedValues: {
      actual_weight: '300',
      estimated_weight: '500',
    },
    context: [
      {
        name: 'Domestic uncrating',
        shipment_type: 'HHG',
      },
    ],
    eventName: o.updateMTOServiceItem,
    tableName: t.mto_service_items,
  };
  it('correctly matches the update service item event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay()).toEqual('Updated service item');
    expect(result.getDetailsLabeledDetails(item)).toMatchObject({
      service_item_name: 'Domestic uncrating',
      shipment_type: 'HHG',
      actual_weight: '300',
      estimated_weight: '500',
    });
  });
});

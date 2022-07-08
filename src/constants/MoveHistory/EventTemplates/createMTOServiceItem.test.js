import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import e from 'constants/MoveHistory/EventTemplates/createMTOServiceItem';

describe('when given a Create basic service item history record', () => {
  const item = {
    action: a.INSERT,
    changedValues: {
      reason: 'Test',
      status: 'SUBMITTED',
    },
    context: [
      {
        name: 'Domestic uncrating',
        shipment_type: 'HHG',
      },
    ],
    eventName: o.createMTOServiceItem,
    tableName: t.mto_service_items,
  };
  it('correctly matches the create service item event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay()).toEqual('Requested service item');
    expect(result.getDetailsLabeledDetails(item)).toMatchObject({
      service_item_name: 'Domestic uncrating',
      shipment_type: 'HHG',
      reason: 'Test',
      status: 'SUBMITTED',
    });
  });
});

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import e from 'constants/MoveHistory/EventTemplates/createMTOServiceItemCustomerContacts';

describe('when given a Create basic service item customer contacts history record', () => {
  const item = {
    action: a.INSERT,
    changedValues: {
      first_available_delivery_date: '2022-06-30T00:00:00+00:00',
      time_military: '1500Z',
      type: 'SECOND',
    },
    eventName: o.createMTOServiceItem,
    tableName: t.mto_service_item_customer_contacts,
  };
  it('correctly matches the create service item customer contacts event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay()).toEqual('Requested service item');
    expect(result.getDetailsLabeledDetails(item)).toMatchObject({
      first_available_delivery_date: '2022-06-30T00:00:00+00:00',
      second_available_delivery_time: '1500Z',
    });
  });
});

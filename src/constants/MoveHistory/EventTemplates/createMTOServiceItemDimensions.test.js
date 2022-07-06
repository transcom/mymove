import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import e from 'constants/MoveHistory/EventTemplates/createMTOServiceItemDimensions';

describe('when given a Create basic service item dimensions history record', () => {
  const item = {
    action: a.INSERT,
    changedValues: {
      height_thousandth_inches: 1000,
      length_thousandth_inches: 3000,
      width_thousandth_inches: 2000,
      type: 'CRATE',
    },
    eventName: o.createMTOServiceItem,
    tableName: t.mto_service_item_dimensions,
  };
  it('correctly matches the create service item dimensions event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay()).toEqual('Requested service item');
    expect(result.getDetailsLabeledDetails(item)).toMatchObject({
      height_thousandth_inches: 1000,
      width_thousandth_inches: 2000,
      length_thousandth_inches: 3000,
      type: 'CRATE',
      crate_size: '1x3x2 in',
    });
  });
});

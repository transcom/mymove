import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import e from 'constants/MoveHistory/EventTemplates/CreateMTOServiceItem/createMTOServiceItemDimensions';

describe('when given a Create basic service item dimensions history record', () => {
  const historyRecord = {
    action: a.INSERT,
    changedValues: {
      height_thousandth_inches: 1000,
      length_thousandth_inches: 3000,
      width_thousandth_inches: 2000,
      type: 'CRATE',
    },
    context: [
      {
        name: 'Domestic uncrating',
        shipment_type: 'HHG',
        shipment_id_abbr: 'a1b2c',
      },
    ],
    eventName: o.createMTOServiceItem,
    tableName: t.mto_service_item_dimensions,
  };
  const template = getTemplate(historyRecord);
  it('correctly matches the create service item dimensions event', () => {
    expect(template).toMatchObject(e);
    expect(template.getEventNameDisplay()).toEqual('Requested service item');
  });
  describe('when given a specific set of details', () => {
    it.each([['crate_size', '1x3x2 in']])(
      'for label %s it displays the proper details value %s',
      async (label, value) => {
        render(template.getDetails(historyRecord));
        expect(screen.getByText(value, { exact: false })).toBeInTheDocument();
      },
    );
  });
});

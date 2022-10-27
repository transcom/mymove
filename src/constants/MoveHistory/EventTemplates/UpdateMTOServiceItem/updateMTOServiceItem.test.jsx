import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import e from 'constants/MoveHistory/EventTemplates/UpdateMTOServiceItem/updateMTOServiceItem';

describe('when given a Update basic service item history record', () => {
  const historyRecord = {
    action: a.UPDATE,
    changedValues: {
      actual_weight: '300',
      estimated_weight: '500',
    },
    context: [
      {
        name: 'Domestic uncrating',
        shipment_type: 'HHG',
        shipment_id_abbr: 'a1b2c',
      },
    ],
    eventName: o.updateMTOServiceItem,
    tableName: t.mto_service_items,
  };
  it('correctly matches the update service item event', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
    expect(template.getEventNameDisplay()).toEqual('Updated service item');
  });
  it('displays shipment type, shipment ID, and service item name properly', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('HHG shipment #A1B2C, Domestic uncrating')).toBeInTheDocument();
  });

  describe('When given a specific set of details', () => {
    it.each([
      ['Actual weight', ': 300'],
      ['Estimated weight', ': 500'],
    ])('displays the proper details value for %s', async (label, value) => {
      const result = getTemplate(historyRecord);
      render(result.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

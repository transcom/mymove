import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateReweigh/updateReweighWeight';

describe('when given an updated reweigh weight', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: { weight: '9001' },
    context: [{ shipment_type: 'HHG', shipment_id_abbr: 'a1b2c' }],
    eventName: 'updateReweigh',
    tableName: 'reweighs',
  };

  it('correctly matches the update reweigh weight event to the proper template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('displays the correct shipment title with ID', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('HHG shipment #A1B2C')).toBeInTheDocument();
  });

  describe('displays the correct labeled values in the details column', () => {
    it.each([['Reweigh weight', ': 9,001 lbs']])('displays the proper details value for %s', async (label, value) => {
      const template = getTemplate(historyRecord);

      render(template.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

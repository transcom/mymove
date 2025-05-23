import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import removeSITExtensionServiceItem from 'constants/MoveHistory/EventTemplates/RemoveSITExtension/removeSITExtensionServiceItem';
import Actions from 'constants/MoveHistory/Database/Actions';

describe('when given a removeSITExtensionServiceItem history record', () => {
  const historyRecord = {
    action: Actions.UPDATE,
    changedValues: {
      status: 'REMOVED',
    },
    context: [
      {
        name: "Domestic origin add'l SIT",
        shipment_id_abbr: '3118a',
        shipment_locator: 'PHD33D-01',
        shipment_type: 'HHG',
      },
    ],
    eventName: 'updateMTOServiceItem',
    tableName: 'sit_extensions',
  };

  it('matches the template from getTemplate', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(removeSITExtensionServiceItem);
  });

  it('returns the correct event display name', () => {
    expect(removeSITExtensionServiceItem.getEventNameDisplay()).toEqual('SIT extension removed');
  });

  it('renders the details via LabeledDetails with merged changed values', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));

    // Check for the presence of the changed values.
    // The actual keys and values displayed depend on your LabeledDetails implementation.
    // Here we expect the values from changedValues to be rendered.
    expect(screen.getByText(/Status/i)).toBeInTheDocument();
    expect(screen.getByText(/REMOVED/i)).toBeInTheDocument();
    expect(screen.getByText('HHG shipment #PHD33D-01')).toBeInTheDocument();
  });
});

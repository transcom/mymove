import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import denySITExtensionServiceItem from 'constants/MoveHistory/EventTemplates/DenySITExtension/denySITExtensionServiceItem';
import Actions from 'constants/MoveHistory/Database/Actions';

describe('when given a denySITExtensionServiceItem history record', () => {
  const historyRecord = {
    action: Actions.UPDATE,
    changedValues: {
      customer_expense: true,
      customer_expense_reason: 'converting to customer expense',
    },
    context: [{ name: 'Domestic destination 1st day SIT' }],
    eventName: 'denySITExtension',
    tableName: 'mto_service_items',
  };

  it('matches the template from getTemplate', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(denySITExtensionServiceItem);
  });

  it('returns the correct event display name', () => {
    expect(denySITExtensionServiceItem.getEventNameDisplay()).toEqual('SIT extension denied');
  });

  it('renders the details via LabeledDetails with merged changed values', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));

    // Check for the presence of the changed values.
    // The actual keys and values displayed depend on your LabeledDetails implementation.
    // Here we expect the values from changedValues to be rendered.
    expect(screen.getByText(/Reason/i)).toBeInTheDocument();
    expect(screen.getByText(/converting to customer expense/i)).toBeInTheDocument();
    // Optionally, check for the boolean value (assuming it renders as a string "true")
    expect(screen.getByText(/Yes/)).toBeInTheDocument();
  });
});

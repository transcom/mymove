import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateBillableWeight/updateBillableWeightAsTIO';

describe('when updating billable weight as a TIO ', () => {
  const historyRecord = {
    action: 'UPDATE',
    // To-do: Update changedValue key to max_billable_weight once available in database.
    changedValues: { authorized_weight: '7999' },
    eventName: 'updateMaxBillableWeightAsTIO',
    tableName: 'entitlements',
  };
  it('correctly matches the update billable weights event', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay(historyRecord)).toEqual('Updated move');
  });
  it('correctly renders the details label component', () => {
    const result = getTemplate(historyRecord);
    render(result.getDetails(historyRecord));
    expect(screen.getByText('Max billable weight')).toBeInTheDocument();
    expect(screen.getByText(': 7,999 lbs')).toBeInTheDocument();
  });
});

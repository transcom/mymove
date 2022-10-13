import { render, screen } from '@testing-library/react';

import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/updateOrders/updateOrder';

describe('when given an Order update history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updateOrder',
    tableName: 'orders',
    detailsType: d.LABELED,
    changedValues: { old_duty_location_id: 'ID1', new_duty_location_id: 'ID2', has_dependents: 'false' },
    context: [{ old_duty_location_name: 'old name', new_duty_location_name: 'new name' }],
  };
  it('correctly matches the Update orders event', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(e);
    render(result.getDetails(historyRecord));
    // expect to have merged context and changedValues
    expect(screen.queryByText('ID1')).not.toBeInTheDocument();
    expect(screen.queryByText('ID2')).not.toBeInTheDocument();
    expect(screen.queryByText('old name')).not.toBeInTheDocument();
    expect(screen.getByText('New duty location name')).toBeInTheDocument();
    expect(screen.getByText('new name', { exact: false })).toBeInTheDocument();
    expect(screen.getByText('Dependents included')).toBeInTheDocument();
  });
});

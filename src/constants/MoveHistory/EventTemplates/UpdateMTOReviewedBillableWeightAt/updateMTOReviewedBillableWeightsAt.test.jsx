import { render, screen } from '@testing-library/react';

import updateMTOReviewedBillableWeightsAt from './updateMTOReviewedBillableWeightsAt';

import getTemplate from 'constants/MoveHistory/TemplateManager';

describe('when given an MTO Reviewed Billable Weight At event', () => {
  const item = {
    action: 'UPDATE',
    eventName: 'updateMTOReviewedBillableWeightsAt',
    tableName: 'moves',
  };
  it('correctly matches the MTO Reviewed Billable Weight At event', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(updateMTOReviewedBillableWeightsAt);
    expect(result.getEventNameDisplay()).toEqual('Updated move');
  });
  it('correctly displays the details component', () => {
    const result = getTemplate(item);
    render(result.getDetails(item));
    expect(screen.getByText('Reviewed weights')).toBeInTheDocument();
  });
});

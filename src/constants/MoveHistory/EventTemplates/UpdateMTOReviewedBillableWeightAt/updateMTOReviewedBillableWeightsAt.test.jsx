import { render, screen } from '@testing-library/react';

import updateMTOReviewedBillableWeightsAt from './updateMTOReviewedBillableWeightsAt';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

describe('when given an MTO Reviewed Billable Weight At event', () => {
  const item = {
    action: a.UPDATE,
    eventName: o.updateMTOReviewedBillableWeightsAt,
    tableName: t.moves,
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

import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import e from 'constants/MoveHistory/EventTemplates/CreateOrders/createOrders';

describe('When given a created orders event for the orders table', () => {
  const item = {
    action: 'INSERT',
    eventName: o.createOrders,
    tableName: t.orders,
    eventNameDisplay: 'Created orders',
    changedValues: {
      status: 'DRAFT',
      report_by_date: '2022-10-18',
      issue_date: '2022-10-11',
      orders_type: 'PERMANENT_CHANGE_OF_STATION',
      origin_duty_location_name: 'Los Angeles AFB',
      new_duty_location_name: 'Fairchild AFB',
      has_dependents: true,
      grade: 'E_1',
    },
    context: [
      {
        new_duty_location_name: 'Fairchild AFB',
        origin_duty_location_name: 'Los Angeles AFB',
      },
    ],
  };
  it('correctly matches to the proper template', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
  });
  describe('When given a specific set of details for created orders', () => {
    it.each([
      ['Status', ': DRAFT'],
      ['Report by date', ': 18 Oct 2022'],
      ['Orders date', ': 11 Oct 2022'],
      ['Orders type', ': Permanent Change Of Station (PCS)'],
      ['Origin duty location name', ': Los Angeles AFB'],
      ['New duty location name', ': Fairchild AFB'],
      ['Dependents included', ': Yes'],
      ['Rank', ': E-1'],
    ])('displays the proper details value for %s', async (label, value) => {
      const result = getTemplate(item);
      render(result.getDetails(item));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

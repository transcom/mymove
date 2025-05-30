import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/CreateOrder/createOrderCreateMoves';

describe('When given a create order event for the moves table from the office side', () => {
  const item = {
    action: 'INSERT',
    eventName: 'createOrder',
    tableName: 'moves',
    eventNameDisplay: 'Created move',
    changedValues: {
      status: 'DRAFT',
    },
    context: [
      {
        counseling_office_name: 'Scott AFB',
      },
    ],
  };
  it('correctly matches to the proper template', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
  });
  describe('When given a specific set of details for created move', () => {
    it.each([
      ['Status', ': DRAFT'],
      ['Counseling office', ': Scott AFB'],
    ])('displays the proper details value for %s', async (label, value) => {
      const result = getTemplate(item);
      render(result.getDetails(item));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

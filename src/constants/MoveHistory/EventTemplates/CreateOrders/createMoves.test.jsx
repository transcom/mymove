import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import e from 'constants/MoveHistory/EventTemplates/CreateOrders/createMoves';

describe('When given a created orders event for the moves table', () => {
  const item = {
    action: 'INSERT',
    eventName: o.createOrders,
    tableName: t.moves,
    eventNameDisplay: 'Created move',
    changedValues: {
      status: 'DRAFT',
    },
  };
  it('correctly matches to the proper template', () => {
    const result = getTemplate(item);
    expect(result).toMatchObject(e);
  });
  describe('When given a specific set of details for created moves', () => {
    it.each([['Status', ': DRAFT']])('displays the proper details value for %s', async (label, value) => {
      const result = getTemplate(item);
      render(result.getDetails(item));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});

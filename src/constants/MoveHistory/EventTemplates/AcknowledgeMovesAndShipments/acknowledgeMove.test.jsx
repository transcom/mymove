import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';

const historyRecord = {
  ACKNOWLEDGE_MOVE: {
    action: a.UPDATE,
    eventName: o.acknowledgeMovesAndShipments,
    tableName: t.moves,
    changedValues: {
      prime_acknowledged_at: '2025-04-13T14:15:33.12345+00:00',
    },
  },
};

describe('When a move is acknowledged by the prime', () => {
  it('displays the prime acknowledged at timestamp', () => {
    const template = getTemplate(historyRecord.ACKNOWLEDGE_MOVE);
    render(template.getDetails(historyRecord.ACKNOWLEDGE_MOVE));
    const label = screen.getByText('Prime Acknowledged At:');
    expect(label).toBeInTheDocument();
    const dateElement = screen.getByText('2025-04-13T14:15:33.12345+00:00');
    expect(dateElement).toBeInTheDocument();
  });
});

/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import EvaluationViolations from './EvaluationViolations';

import { MockProviders } from 'testUtils';
// import { useEvaluationViolationsQueries } from 'hooks/queries';

// jest.mock('hooks/queries', () => ({
//   useEvaluationViolationsQueries: jest.fn(),
// }));

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  // useLocation: () => ({
  //   pathname: `/moves/${mockRequestedMoveCode}/evaluation-reports`,
  //   state: { showDeleteSuccess: true },
  // }),
  useParams: () => {
    return { moveCode: 'ABCDEFGH', reportId: 'db30c135-1d6d-4a0d-a6d5-f408474f6ee2' };
  },
}));

// TODO: Figure out mock report data.
describe('EvaluationViolations', () => {
  it('renders the component content', async () => {
    render(
      <MockProviders
        initialEntries={['/moves/ABCDEFGH/evaluation-reports/db30c135-1d6d-4a0d-a6d5-f408474f6ee2/violations']}
      >
        <EvaluationViolations />
      </MockProviders>,
    );

    // Check out headings
    expect(await screen.getByRole('heading', { name: 'REPORT ID #', level: 6 })).toBeInTheDocument();
    expect(await screen.getByRole('heading', { name: 'MOVE CODE', level: 6 })).toBeInTheDocument();
    expect(await screen.getByRole('heading', { name: 'MTO REFERENCE ID #', level: 6 })).toBeInTheDocument();

    expect(await screen.getByRole('heading', { name: 'Select violations', level: 2 })).toBeInTheDocument();

    // Check out buttons
    const buttons = await screen.getAllByRole('button');
    expect(buttons).toHaveLength(4);
    expect(buttons[0]).toHaveTextContent('< Back to Evaluation form');
    expect(buttons[1]).toHaveTextContent('Cancel');
    expect(buttons[2]).toHaveTextContent('Save draft');
    expect(buttons[3]).toHaveTextContent('Review and submit');
  });

  // TODO: This test is not done yet
  it('reroutes back to the eval report', async () => {
    render(
      <MockProviders
        initialEntries={['/moves/ABCDEFGH/evaluation-reports/db30c135-1d6d-4a0d-a6d5-f408474f6ee2/violations']}
      >
        <EvaluationViolations />
      </MockProviders>,
    );

    const buttons = await screen.getAllByRole('button');
    userEvent.click(buttons[0]);

    // expect(window.location.pathname).toBe(`/moves/ABCDEFGH/evaluation-reports/db30c135-1d6d-4a0d-a6d5-f408474f6ee2`);

    // screen.debug();
  });
});

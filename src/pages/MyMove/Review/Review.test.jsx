/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ConnectedReview from 'pages/MyMove/Review/Review';
import { MockProviders } from 'testUtils';

// Mock the summary part of the review page since we're just testing the
// navigation portion.
jest.mock('components/Customer/Review/Summary/index', () => 'summary');

describe('Review page', () => {
  const testProps = {
    canMoveNext: true,
    push: jest.fn(),
    match: {
      path: '/moves/:moveId/review',
      url: '/moves/3a8c9f4f-7344-4f18-9ab5-0de3ef57b901/review',
      isExact: true,
      params: {
        moveId: '3a8c9f4f-7344-4f18-9ab5-0de3ef57b901',
      },
    },
  };

  it('renders the Review Page', async () => {
    const { findByRole } = render(
      <MockProviders>
        <ConnectedReview {...testProps} />
      </MockProviders>,
    );

    await findByRole('heading', { level: 1, name: 'Review your details' });
  });

  it('Finish Later button goes to the home page', async () => {
    const { findByRole } = render(
      <MockProviders>
        <ConnectedReview {...testProps} />
      </MockProviders>,
    );

    const backButton = await findByRole('button', { name: 'Finish later' });

    expect(backButton).toBeInTheDocument();

    userEvent.click(backButton);

    expect(testProps.push).toHaveBeenCalledWith('/');
  });

  it('next button goes to the Agreement page', async () => {
    const { findByRole } = render(
      <MockProviders>
        <ConnectedReview {...testProps} />
      </MockProviders>,
    );

    const submitButton = await findByRole('button', { name: 'Next' });

    expect(submitButton).toBeInTheDocument();

    userEvent.click(submitButton);

    expect(testProps.push).toHaveBeenCalledWith(`/moves/${testProps.match.params.moveId}/agreement`);
  });

  afterEach(jest.resetAllMocks);
});

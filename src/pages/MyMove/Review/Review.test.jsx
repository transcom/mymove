import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ConnectedReview from 'pages/MyMove/Review/Review';
import { MockProviders } from 'testUtils';

// Mock the summary part of the review page since we're just testing the
// navigation portion.
jest.mock('components/Customer/Review/Summary/Summary', () => 'summary');

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

  const testState = {
    flash: {
      flashMessage: {
        type: 'SET_FLASH_MESSAGE',
        title: 'Details saved',
        messageType: 'success',
        message: 'Review your info and submit your move request now, or come back and finish later.',
        key: 'PPM_ONBOARDING_SUBMIT_SUCCESS',
        slim: false,
      },
    },
  };

  it('renders the Review Page', async () => {
    render(
      <MockProviders>
        <ConnectedReview {...testProps} />
      </MockProviders>,
    );

    await screen.findByRole('heading', { level: 1, name: 'Review your details' });
  });

  it('Finish Later button goes to the home page', async () => {
    render(
      <MockProviders>
        <ConnectedReview {...testProps} />
      </MockProviders>,
    );

    const backButton = await screen.findByRole('button', { name: 'Finish later' });

    expect(backButton).toBeInTheDocument();

    userEvent.click(backButton);

    expect(testProps.push).toHaveBeenCalledWith('/');
  });

  it('next button goes to the Agreement page', async () => {
    render(
      <MockProviders>
        <ConnectedReview {...testProps} />
      </MockProviders>,
    );

    const submitButton = await screen.findByRole('button', { name: 'Next' });

    expect(submitButton).toBeInTheDocument();

    userEvent.click(submitButton);

    expect(testProps.push).toHaveBeenCalledWith(`/moves/${testProps.match.params.moveId}/agreement`);
  });

  it('renders the success alert flash message', async () => {
    render(
      <MockProviders initialState={testState}>
        <ConnectedReview {...testProps} />
      </MockProviders>,
    );

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('Details saved');
    expect(
      screen.getByText('Review your info and submit your move request now, or come back and finish later.'),
    ).toBeInTheDocument();
  });

  afterEach(jest.resetAllMocks);
});

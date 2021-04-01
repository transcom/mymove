/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ConnectedReview from './index';

import { MockProviders } from 'testUtils';

// Mock the summary part of the review page since we're just testing the
// navigation portion.
jest.mock('components/Customer/Review/Summary/index', () => 'summary');

describe('Review page', () => {
  const testProps = {
    canMoveNext: true,
    push: jest.fn(),
  };

  it('renders the Review Page', async () => {
    await waitFor(() => {
      const wrapper = mount(
        <MockProviders>
          <ConnectedReview {...testProps} />
        </MockProviders>,
      );
      expect(wrapper.find('h1').text()).toBe('Review your details');
    });
  });

  it('Finish Later button goes to the home page', async () => {
    const { queryByText } = render(
      <MockProviders>
        <ConnectedReview {...testProps} />
      </MockProviders>,
    );

    const backButton = queryByText('Finish later');

    await waitFor(() => {
      expect(backButton).toBeInTheDocument();
    });

    userEvent.click(backButton);
    expect(testProps.push).toHaveBeenCalledWith('/');
  });

  it('next button goes to the Agreement page', async () => {
    const { queryByText } = render(
      <MockProviders>
        <ConnectedReview {...testProps} />
      </MockProviders>,
    );

    const submitButton = queryByText('Next');
    expect(submitButton).toBeInTheDocument();
    userEvent.click(submitButton);

    expect(testProps.push).toHaveBeenCalledWith('/moves/:moveId/agreement');
  });

  afterEach(jest.resetAllMocks);
});

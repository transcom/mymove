import React from 'react';
import { render, screen } from '@testing-library/react';

import SomethingWentWrong from './index';

describe('SomethingWentWrong tests', () => {
  it('renders without crashing', async () => {
    const { container } = render(<SomethingWentWrong />);

    const errorPage = await container.querySelector('.usa-grid');
    expect(errorPage).toBeInTheDocument();
  });

  it('should render the correct image on the page', () => {
    render(<SomethingWentWrong />);

    const image = screen.getByRole('img');
    expect(image).toBeInTheDocument();
    expect(image).toHaveAttribute('src', 'sad-computer.png');
  });

  it('should render the correct text on the page', () => {
    render(<SomethingWentWrong />);

    const oopsMsg = screen.getByRole('heading', { level: 2 });
    expect(oopsMsg).toBeInTheDocument();
    expect(oopsMsg).toHaveTextContent('Oops!Something went wrong.');

    const tryAgainMsg = screen.getByText('Please try again in a few moments.');
    expect(tryAgainMsg).toBeInTheDocument();

    const contactMsg = screen.getByTestId('contactMsg');
    expect(contactMsg).toBeInTheDocument();
    expect(contactMsg).toHaveTextContent(
      'If you continue to receive this error, call (800) 462-2176, Option 2 or email us.',
    );

    const email = screen.getByRole('link', { name: 'email us' });
    expect(email).toBeInTheDocument();
    expect(email).toHaveAttribute('href', 'mailto:usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil');
  });
});

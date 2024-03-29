import React from 'react';
import { render, screen } from '@testing-library/react';
import SmartCardRedirect from './SmartCardRedirect';

describe('SmartCardRedirect tests', () => {
  it('renders without crashing', async () => {
    const { container } = render(<SmartCardRedirect />);

    const errorPage = await container.querySelector('.usa-grid');
    expect(errorPage).toBeInTheDocument();
  });

  it('should render the smart card image on the page', () => {
    render(<SmartCardRedirect />);

    const image = screen.getByRole('img');
    expect(image).toBeInTheDocument();
    expect(image).toHaveAttribute('src', 'smart-card.png');
  });

  it('should render the text on the page', () => {
    render(<SmartCardRedirect />);

    const oopsMsg = screen.getByRole('heading', { level: 2 });
    expect(oopsMsg).toBeInTheDocument();
    expect(oopsMsg).toHaveTextContent('You must sign in with your smart card first.');

    const helperText = screen.getByTestId('helperText');
    expect(helperText).toBeInTheDocument();

    const contactMsg = screen.getByTestId('contactMsg');
    expect(contactMsg).toBeInTheDocument();

    const email = screen.getByRole('link', { name: 'email us' });
    expect(email).toBeInTheDocument();
    expect(email).toHaveAttribute('href', 'mailto:usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil');
  });
});

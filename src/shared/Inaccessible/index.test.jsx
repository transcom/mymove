import React from 'react';
import { render, screen } from '@testing-library/react';

import Inaccessible from './index';

describe('Inaccessible tests', () => {
  it('renders without crashing', async () => {
    const { container } = render(<Inaccessible />);

    const errorPage = await container.querySelector('.usa-grid');
    expect(errorPage).toBeInTheDocument();
  });

  it('should render the correct image on the page', () => {
    render(<Inaccessible />);

    const image = screen.getByRole('img');
    expect(image).toBeInTheDocument();
    expect(image).toHaveAttribute('src', 'sad-computer.png');
  });

  it('should render the correct text on the page', () => {
    render(<Inaccessible />);

    const inaccessibleMsg = screen.getByRole('heading', { level: 2 });
    expect(inaccessibleMsg).toBeInTheDocument();
    expect(inaccessibleMsg).toHaveTextContent('Page is not accessible.');

    const contactMsg = screen.getByTestId('contactMsg');
    expect(contactMsg).toBeInTheDocument();
    expect(contactMsg).toHaveTextContent(
      'If you feel this message was received in error, please call (800) 462-2176, Option 2 or email us.',
    );

    const email = screen.getByRole('link', { name: 'email us' });
    expect(email).toBeInTheDocument();
    expect(email).toHaveAttribute('href', 'mailto:usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil');
  });
});

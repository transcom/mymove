import React from 'react';
import { MemoryRouter } from 'react-router-dom';
import { render, screen } from '@testing-library/react';

import ContactInfoDisplay from './ContactInfoDisplay';

describe('ContactInfoDisplay Component', () => {
  const mockUserInfo = {
    name: 'John Doe',
    telephone: '123-456-7890',
    email: 'john.doe@example.com',
  };
  const mockEditURL = '/edit-profile';

  it('renders without crashing', () => {
    render(
      <MemoryRouter>
        <ContactInfoDisplay officeUserInfo={mockUserInfo} editURL={mockEditURL} />
      </MemoryRouter>,
    );

    // Check for the heading
    expect(screen.getByRole('heading', { name: 'Contact info' })).toBeInTheDocument();

    // Check the name, email, and phone values
    expect(screen.getByTestId('name')).toHaveTextContent(mockUserInfo.name);
    expect(screen.getByTestId('email')).toHaveTextContent(mockUserInfo.email);
    expect(screen.getByTestId('phone')).toHaveTextContent(mockUserInfo.telephone);

    // Check the Edit link
    const editLink = screen.getByRole('link', { name: 'Edit' });
    expect(editLink).toBeInTheDocument();
    expect(editLink).toHaveAttribute('href', mockEditURL);
  });

  it('displays the correct data from props', () => {
    render(
      <MemoryRouter>
        <ContactInfoDisplay officeUserInfo={mockUserInfo} editURL={mockEditURL} />
      </MemoryRouter>,
    );

    expect(screen.getByTestId('name')).toHaveTextContent('John Doe');
    expect(screen.getByTestId('email')).toHaveTextContent('john.doe@example.com');
    expect(screen.getByTestId('phone')).toHaveTextContent('123-456-7890');
  });

  it('renders the Edit link with correct URL', () => {
    render(
      <MemoryRouter>
        <ContactInfoDisplay officeUserInfo={mockUserInfo} editURL={mockEditURL} />
      </MemoryRouter>,
    );

    const editLink = screen.getByRole('link', { name: 'Edit' });
    expect(editLink).toHaveAttribute('href', '/edit-profile');
  });
});

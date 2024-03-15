import React from 'react';
import { render, screen } from '@testing-library/react';

import NameForm from 'components/Customer/NameForm/NameForm';

describe('Name page', () => {
  it('renders the NameForm', async () => {
    render(<NameForm />);

    expect(await screen.findByRole('heading', { name: 'Name', level: 1 })).toBeInTheDocument();
  });
});

import React from 'react';
import { screen, waitFor, render } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ServiceOrderNumberModal from './ServiceOrderNumberModal';

describe('components/Office/ServiceOrderNumberModal', () => {
  it('is connected', () => {
    render(<ServiceOrderNumberModal />);
    expect(screen.queryByRole('heading', { name: 'Edit service order number' })).not.toBeInTheDocument();
  });

  it('defaults values from props', () => {
    render(<ServiceOrderNumberModal isOpen serviceOrderNumber="ABC123" />);
    expect(screen.getByRole('heading', { name: 'Edit service order number' })).toBeInTheDocument();
    expect(screen.getByRole('textbox')).toHaveValue('ABC123');
  });

  it('validates input', async () => {
    const onSubmit = jest.fn();
    render(<ServiceOrderNumberModal isOpen onSubmit={onSubmit} />);

    const textbox = screen.getByRole('textbox');
    const saveButton = screen.getByRole('button', { name: 'Save' });

    userEvent.type(textbox, '!');
    userEvent.click(saveButton);
    await waitFor(() => expect(screen.getByText('Letters and numbers only')).toBeInTheDocument());
    expect(onSubmit).not.toHaveBeenCalled();

    userEvent.clear(textbox);
    userEvent.click(saveButton);
    await waitFor(() => expect(screen.getByText('Required')).toBeInTheDocument());
    expect(onSubmit).not.toHaveBeenCalled();

    userEvent.type(textbox, 'ABC123');
    userEvent.click(saveButton);
    await waitFor(() => expect(screen.queryByTestId('errorMessage')).not.toBeInTheDocument());
    expect(onSubmit).toHaveBeenCalledWith({ serviceOrderNumber: 'ABC123' });
  });
});

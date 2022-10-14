import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { DocsButton, EditButton } from './IconButtons';

describe('DocsButton', () => {
  it('should render the button', () => {
    render(<DocsButton label="my docs button" />);
    expect(screen.getByRole('button')).toHaveTextContent('my docs button');
    expect(screen.getByTestId('docs-icon')).toHaveClass('fa-file');
  });
  it('should pass props down', () => {
    render(<DocsButton label="my docs button" className="sample-class" />);
    expect(screen.getByRole('button')).toHaveClass('sample-class');
  });
  it('onClick works', async () => {
    const mockFn = jest.fn();
    render(<DocsButton label="my docs button" onClick={mockFn} />);
    await userEvent.click(screen.getByRole('button'));
    await waitFor(() => {
      expect(mockFn).toHaveBeenCalled();
    });
  });
});

describe('EditButton', () => {
  it('should render the button', () => {
    render(<EditButton label="my edit button" />);
    expect(screen.getByRole('button')).toHaveTextContent('my edit button');
    expect(screen.getByTestId('edit-icon')).toHaveClass('fa-pen');
  });
  it('should pass props down', () => {
    render(<EditButton label="my edit button" className="sample-class" />);
    expect(screen.getByRole('button')).toHaveClass('sample-class');
  });
  it('onClick works', async () => {
    const mockFn = jest.fn();
    render(<EditButton label="my edit button" onClick={mockFn} />);
    await userEvent.click(screen.getByRole('button'));
    await waitFor(() => {
      expect(mockFn).toHaveBeenCalled();
    });
  });
});

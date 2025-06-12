import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import YesNoBoolean from './YesNoBoolean';
import { renderWithProviders } from 'testUtils';

describe('YesNoBoolean Component', () => {
  const defaultProps = {
    onChange: jest.fn(),
    value: false,
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders yes and no radio buttons', () => {
    renderWithProviders(<YesNoBoolean {...defaultProps} />);

    expect(screen.getByLabelText('Yes')).toBeInTheDocument();
    expect(screen.getByLabelText('No')).toBeInTheDocument();
  });

  it('should reflect initial value prop correctly when true', () => {
    renderWithProviders(<YesNoBoolean {...defaultProps} value={true} />);

    expect(screen.getByLabelText('Yes')).toBeChecked();
    expect(screen.getByLabelText('No')).not.toBeChecked();
  });

  it('should reflects initial value prop correctly when false', () => {
    renderWithProviders(<YesNoBoolean {...defaultProps} value={false} />);

    expect(screen.getByLabelText('Yes')).not.toBeChecked();
    expect(screen.getByLabelText('No')).toBeChecked();
  });

  it('handles onChange when selecting Yes', () => {
    renderWithProviders(<YesNoBoolean {...defaultProps} value={false} />);

    const yesRadio = screen.getByLabelText('Yes');
    fireEvent.click(yesRadio);

    expect(defaultProps.onChange).toHaveBeenCalledWith(true);
    expect(defaultProps.onChange).toHaveBeenCalledTimes(1);
  });

  it('handles onChange when selecting No', () => {
    render(<YesNoBoolean {...defaultProps} value={true} />);

    const noRadio = screen.getByLabelText('No');
    fireEvent.click(noRadio);

    expect(defaultProps.onChange).toHaveBeenCalledWith(false);
    expect(defaultProps.onChange).toHaveBeenCalledTimes(1);
  });

  it('works with input prop object', () => {
    const inputProps = {
      input: {
        value: 'true', // String value should be converted to boolean
        onChange: jest.fn(),
      },
    };

    renderWithProviders(<YesNoBoolean {...inputProps} />);

    expect(screen.getByLabelText('Yes')).toBeChecked();
    expect(screen.getByLabelText('No')).not.toBeChecked();

    fireEvent.click(screen.getByLabelText('No'));
    expect(inputProps.input.onChange).toHaveBeenCalledWith(false);
  });

  it('applies correct classes to inputs and labels', () => {
    renderWithProviders(<YesNoBoolean {...defaultProps} />);

    const radioInputs = screen.getAllByRole('radio');
    const labels = screen.getAllByText(/Yes|No/);

    radioInputs.forEach((input) => {
      expect(input).toHaveClass('usa-radio__input', 'inline_radio');
    });

    labels.forEach((label) => {
      expect(label).toHaveClass('usa-radio__label', 'inline_radio');
    });
  });

  it('has unique IDs for each radio button', () => {
    renderWithProviders(<YesNoBoolean {...defaultProps} />);
    const yesRadio = screen.getByLabelText('Yes');
    const noRadio = screen.getByLabelText('No');

    const yesId = yesRadio.id;
    const noId = noRadio.id;

    expect(yesId).not.toBe(noId);
    expect(yesId).toMatch(/^yes_no_/);
    expect(noId).toMatch(/^yes_no_/);
  });
});

import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import RadioButton from '.';
import { renderWithProviders } from 'testUtils';

describe('RadioButton Component', () => {
  const defaultProps = {
    name: 'test-radio-btn',
    label: 'Test Label',
    onChange: jest.fn(),
    value: 'test-value',
    checked: false,
  };

  // Reset mocks before each test
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders radio button with label', () => {
    renderWithProviders(<RadioButton {...defaultProps} />);

    const radioInput = screen.getByRole('radio');
    const label = screen.getByText('Test Label');

    expect(radioInput).toBeInTheDocument();
    expect(label).toBeInTheDocument();
  });

  it('applies correct props to input element', () => {
    renderWithProviders(<RadioButton {...defaultProps} />);

    const radioInput = screen.getByRole('radio');

    expect(radioInput).toHaveAttribute('name', 'test-radio-btn');
    expect(radioInput).toHaveAttribute('value', 'test-value');
    expect(radioInput).not.toBeChecked();
  });

  it('reflects checked state when true', () => {
    renderWithProviders(<RadioButton {...defaultProps} checked={true} />);

    const radioInput = screen.getByRole('radio');
    expect(radioInput).toBeChecked();
  });

  it('calls onChange handler when clicked', () => {
    renderWithProviders(<RadioButton {...defaultProps} />);

    const radioInput = screen.getByRole('radio');
    fireEvent.click(radioInput);

    expect(defaultProps.onChange).toHaveBeenCalledTimes(1);
  });

  it('applies custom input className', () => {
    const inputClassName = 'custom-input-class';
    renderWithProviders(<RadioButton {...defaultProps} inputClassName={inputClassName} />);

    const radioInput = screen.getByRole('radio');
    expect(radioInput).toHaveClass(inputClassName);
  });

  it('applies custom label className', () => {
    const labelClassName = 'custom-label-class';
    renderWithProviders(<RadioButton {...defaultProps} labelClassName={labelClassName} />);

    const label = screen.getByText('Test Label');
    expect(label).toHaveClass(labelClassName);
  });

  it('uses provided testId for testing', () => {
    const testId = 'custom-test-id';
    render(<RadioButton {...defaultProps} testId={testId} />);

    const radioInput = screen.getByTestId(testId);
    expect(radioInput).toBeInTheDocument();
  });

  it('associates label with input via matching ids', () => {
    renderWithProviders(<RadioButton {...defaultProps} />);

    const radioInput = screen.getByRole('radio');
    const label = screen.getByText('Test Label');

    const inputId = radioInput.getAttribute('id');
    const labelFor = label.getAttribute('for');

    expect(inputId).toBe(labelFor);
    expect(inputId).toBeTruthy();
  });
});

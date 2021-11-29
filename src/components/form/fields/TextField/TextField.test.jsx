import React from 'react';
import { render } from '@testing-library/react';
import { useField } from 'formik'; // package will be auto mocked

import TextField from './TextField';

// mock out formik hook as we are not testing formik
// needs to be before first describe
jest.mock('formik');

describe('TextField component', () => {
  it('renders the elements that make up a field', () => {
    const mockMeta = {
      touched: false,
      error: '',
      initialError: '',
      initialTouched: false,
      initialValue: '',
      value: '',
    };
    const mockField = {
      value: '',
      checked: false,
      onChange: jest.fn(),
      onBlur: jest.fn(),
      multiple: undefined,
      name: 'firstName',
    };

    useField.mockReturnValue([mockField, mockMeta]);

    const { queryByText, queryByLabelText } = render(
      <TextField name="firstName" label="First Name" type="text" id="firstName" />,
    );

    expect(queryByText('First Name')).toBeInstanceOf(HTMLLabelElement);
    expect(queryByLabelText('First Name')).toBeInstanceOf(HTMLInputElement);
    expect(queryByLabelText('First Name')).toHaveAttribute('name', 'firstName');
    expect(queryByLabelText('First Name')).toHaveAttribute('id', 'firstName');
  });

  it('passes a custom className prop to the input element', () => {
    useField.mockReturnValue([{}, {}]);

    const { queryByLabelText } = render(
      <TextField name="firstName" className="myCustomInputClass" label="First Name" type="text" id="firstName" />,
    );

    expect(queryByLabelText('First Name')).toHaveClass('myCustomInputClass');
  });

  describe('with an error message', () => {
    it('does not show the error message if the input is untouched', () => {
      const mockMeta = {
        touched: false,
        error: 'This field is required',
        initialError: '',
        initialTouched: false,
        initialValue: '',
        value: '',
      };

      const mockField = {
        value: '',
        checked: false,
        onChange: jest.fn(),
        onBlur: jest.fn(),
        multiple: undefined,
        name: 'firstName',
      };

      useField.mockReturnValue([mockField, mockMeta]);

      const { queryByText } = render(<TextField name="firstName" label="First Name" type="text" id="firstName" />);
      expect(queryByText('First Name')).not.toHaveClass('usa-label--error');
      expect(queryByText('This field is required')).not.toBeInTheDocument();
    });

    it('shows the error message if the input is touched', () => {
      const mockMeta = {
        touched: true,
        error: 'This field is required',
        initialError: '',
        initialTouched: false,
        initialValue: '',
        value: '',
      };

      const mockField = {
        value: '',
        checked: false,
        onChange: jest.fn(),
        onBlur: jest.fn(),
        multiple: undefined,
        name: 'firstName',
      };

      useField.mockReturnValue([mockField, mockMeta]);

      const { queryByText } = render(<TextField name="firstName" label="First Name" type="text" id="firstName" />);

      expect(queryByText('First Name')).toHaveClass('usa-label--error');
      expect(queryByText('This field is required')).toBeInTheDocument();
    });
  });

  describe('with a warning', () => {
    it('shows the warning if there is no error shown', () => {
      const mockMeta = {
        touched: false,
        error: 'This field is required',
        initialError: '',
        initialTouched: false,
        initialValue: '',
        value: '',
      };

      const mockField = {
        value: '',
        checked: false,
        onChange: jest.fn(),
        onBlur: jest.fn(),
        multiple: undefined,
        name: 'firstName',
      };

      useField.mockReturnValue([mockField, mockMeta]);

      const { queryByText } = render(
        <TextField name="firstName" label="First Name" type="text" id="firstName" warning="This is a warning" />,
      );

      expect(queryByText('This is a warning')).toBeInTheDocument();
    });

    it('does not show the warning if an error is showing', () => {
      const mockMeta = {
        touched: true,
        error: 'This field is required',
        initialError: '',
        initialTouched: false,
        initialValue: '',
        value: '',
      };

      const mockField = {
        value: '',
        checked: false,
        onChange: jest.fn(),
        onBlur: jest.fn(),
        multiple: undefined,
        name: 'firstName',
      };

      useField.mockReturnValue([mockField, mockMeta]);

      const { queryByText } = render(
        <TextField name="firstName" label="First Name" type="text" id="firstName" warning="This is a warning" />,
      );

      expect(queryByText('This is a warning')).not.toBeInTheDocument();
    });
  });

  afterEach(jest.resetAllMocks);
});

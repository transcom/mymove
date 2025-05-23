import React from 'react';
import { render } from '@testing-library/react';
import { useField } from 'formik'; // package will be auto mocked
import userEvent from '@testing-library/user-event';

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

    expect(queryByText('First Name').parentElement).toBeInstanceOf(HTMLLabelElement);
    expect(queryByLabelText('First Name')).toBeInstanceOf(HTMLInputElement);
    expect(queryByLabelText('First Name')).toHaveAttribute('name', 'firstName');
    expect(queryByLabelText('First Name')).toHaveAttribute('id', 'firstName');
  });

  it('renders the required red asterisk when prop is provided', () => {
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

    const { getByTestId } = render(
      <TextField name="firstName" label="First Name" type="text" id="firstName" required showRequiredAsterisk />,
    );
    expect(getByTestId('requiredAsterisk')).toBeInTheDocument();
  });

  it('passes a custom className prop to the input element', () => {
    useField.mockReturnValue([{}, {}]);

    const { queryByLabelText } = render(
      <TextField name="firstName" className="myCustomInputClass" label="First Name" type="text" id="firstName" />,
    );

    expect(queryByLabelText('First Name')).toHaveClass('myCustomInputClass');
  });

  it('can include a trailing button', async () => {
    useField.mockReturnValue([{}, {}]);

    const mockOnButtonClick = jest.fn();
    const { getByTestId } = render(
      <TextField
        name="testName"
        id="testId"
        label="testLabel"
        button={
          <button type="button" data-testid="testButton" onClick={mockOnButtonClick}>
            Test button
          </button>
        }
      />,
    );

    // Verify button is shown
    expect(getByTestId('testButton')).toBeInTheDocument();
    expect(getByTestId('testButton')).toHaveTextContent('Test button');

    // Click the button
    expect(mockOnButtonClick).not.toHaveBeenCalled();
    await userEvent.click(getByTestId('testButton'));
    expect(mockOnButtonClick).toHaveBeenCalled();
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

      expect(queryByText('First Name').parentElement).toHaveClass('usa-label--error');
      expect(queryByText('This field is required')).toBeInTheDocument();
    });

    it('renders a prefix before the input field', () => {
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
        name: 'prefixedInput',
      };

      useField.mockReturnValue([mockField, mockMeta]);

      const { getByText, getByLabelText } = render(
        <TextField
          name="prefixedInput"
          label="Prefixed Input"
          type="text"
          id="prefixedInput"
          prefix="TERMINATED FOR CAUSE:"
        />,
      );

      // Check the prefix span is rendered
      expect(getByText('TERMINATED FOR CAUSE:')).toBeInTheDocument();
      expect(getByLabelText('Prefixed Input')).toBeInstanceOf(HTMLInputElement);
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

  describe('as a text area', () => {
    it('is of type HTMLTextArea', () => {
      useField.mockReturnValue([{}, {}]);

      const { queryByLabelText } = render(
        <TextField
          name="firstName"
          className="myCustomInputClass"
          label="First Name"
          type="text"
          id="firstName"
          display="textarea"
        />,
      );

      expect(queryByLabelText('First Name')).toBeInstanceOf(HTMLTextAreaElement);
    });
  });

  describe('as a read only', () => {
    it('is of type HTMLLabelElement', () => {
      useField.mockReturnValue([{}, {}]);

      const { queryByTestId } = render(
        <TextField
          name="firstName"
          className="myCustomInputClass"
          label="First Name"
          type="text"
          id="firstName"
          display="readonly"
        />,
      );

      expect(queryByTestId('First Name')).toBeInstanceOf(HTMLLabelElement);
    });
  });

  afterEach(jest.resetAllMocks);
});

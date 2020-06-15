import React from 'react';
import { shallow } from 'enzyme';
import { FormGroup, Label, TextInput as UswdsTextInput } from '@trussworks/react-uswds';

import { ErrorMessage } from '../index';

import { TextInput, TextInputMinimal } from './index';

const mockOnChange = jest.fn();
// mock out formik hook as we are not testing formik
// needs to be before first describe
jest.mock('formik', () => {
  return {
    ...jest.requireActual('formik'),
    useField: () => [
      {
        onChange: mockOnChange,
      },
      { touched: true, error: 'sample error' },
    ],
  };
});

describe('TextInputMinimal', () => {
  describe('with name prop', () => {
    const wrapper = shallow(<TextInputMinimal className="sample-class" name="firstName" type="text" id="firstName" />);

    it('should render an ErrorMessage', () => {
      const errorMessage = wrapper.find(ErrorMessage);
      expect(errorMessage.length).toBe(1);
      expect(errorMessage.prop('display')).toBe(true);
      expect(errorMessage.prop('children')).toBe('sample error');
    });

    it('should render a USWDS TextInput', () => {
      const textInput = wrapper.find(UswdsTextInput);
      expect(textInput.length).toBe(1);
      expect(textInput.prop('className')).toBe('sample-class');
      expect(textInput.prop('type')).toBe('text');
    });

    it('should trigger onChange properly', () => {
      const textInput = wrapper.find(UswdsTextInput);
      expect(textInput.prop('onChange')).toBe(mockOnChange);
      textInput.simulate('change', { value: 'sample' });
      expect(mockOnChange).toHaveBeenCalledWith({ value: 'sample' });
    });
  });

  describe('with id prop', () => {
    const wrapper = shallow(<TextInputMinimal className="sample-class" id="lastName" type="text" name="lastName" />);

    it('should render an ErrorMessage', () => {
      const errorMessage = wrapper.find(ErrorMessage);
      expect(errorMessage.length).toBe(1);
      expect(errorMessage.prop('display')).toBe(true);
      expect(errorMessage.prop('children')).toBe('sample error');
    });

    it('should render a USWDS TextInput', () => {
      const textInput = wrapper.find(UswdsTextInput);
      expect(textInput.length).toBe(1);
      expect(textInput.prop('id')).toBe('lastName');
    });
  });

  describe('with no id or name prop', () => {
    it('should render console error', () => {
      const spy = jest.spyOn(global.console, 'error');
      shallow(<TextInputMinimal className="sample-class" type="text" />);

      expect(spy).toHaveBeenCalledWith(
        expect.stringMatching(
          /The prop `id` is marked as required in `TextInputMinimal`, but its value is `undefined`/,
        ),
      );
      expect(spy).toHaveBeenCalledWith(
        expect.stringMatching(
          /The prop `name` is marked as required in `TextInputMinimal`, but its value is `undefined`/,
        ),
      );
    });
  });

  afterEach(jest.resetAllMocks);
});

describe('TextInput', () => {
  describe('with name prop', () => {
    const wrapper = shallow(
      <TextInput className="sample-class" name="firstName" label="First Name" type="text" id="firstName" />,
    );

    it('should render a FormGroup', () => {
      const group = wrapper.find(FormGroup);
      expect(group.length).toBe(1);
      expect(group.prop('error')).toBe(true);
    });

    it('should render a Label', () => {
      const label = wrapper.find(FormGroup).find(Label);
      expect(label.length).toBe(1);
      expect(label.prop('error')).toBe(true);
      expect(label.prop('htmlFor')).toBe('firstName');
      expect(label.prop('children')).toBe('First Name');
    });

    it('should render a TextInputMinimal', () => {
      const textInputMinimal = wrapper.find(FormGroup).find(TextInputMinimal);
      expect(textInputMinimal.length).toBe(1);
      expect(textInputMinimal.prop('name')).toBe('firstName');
      expect(textInputMinimal.prop('type')).toBe('text');
      expect(textInputMinimal.prop('className')).toBe('sample-class');
    });
  });

  describe('with id prop', () => {
    const wrapper = shallow(
      <TextInput className="sample-class" id="lastName" label="Last Name" type="text" name="lastName" />,
    );

    it('should render a Label', () => {
      const label = wrapper.find(FormGroup).find(Label);
      expect(label.length).toBe(1);
      expect(label.prop('htmlFor')).toBe('lastName');
    });

    it('should render a TextInputMinimal', () => {
      const textInput = wrapper.find(FormGroup).find(TextInputMinimal);
      expect(textInput.length).toBe(1);
      expect(textInput.prop('id')).toBe('lastName');
    });
  });

  describe('with no id or name prop', () => {
    it('should render console error', () => {
      const spy = jest.spyOn(global.console, 'error');
      shallow(<TextInput className="sample-class" label="Some Name" type="text" />);

      expect(spy).toHaveBeenCalledWith(
        expect.stringMatching(/The prop `id` is marked as required in `TextInput`, but its value is `undefined`/),
      );
      expect(spy).toHaveBeenCalledWith(
        expect.stringMatching(/The prop `name` is marked as required in `TextInput`, but its value is `undefined`/),
      );
    });
  });

  afterEach(jest.resetAllMocks);
});

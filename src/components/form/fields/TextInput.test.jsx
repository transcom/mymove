import React from 'react';
import { shallow } from 'enzyme';
import { FormGroup, Label, TextInput as UswdsTextInput } from '@trussworks/react-uswds';
import { ErrorMessage } from '..';

describe('TextInput', () => {
  // mock out formik hook as we are not testing formik
  // needs to be before first describe
  jest.mock('formik', () => {
    return {
      useField: jest.fn().mockReturnValue([
        {
          onChange: jest.fn().mockName('onChange'),
        },
        { touched: true, error: 'sample error' },
      ]),
    };
  });
  // require the above mock for expectations
  // eslint-disable-next-line global-require
  const mock = require('formik');

  // import component we are testing after mock created
  // eslint-disable-next-line global-require
  const { TextInput } = require('.');

  it('should call useField', () => {
    expect(mock.useField).toHaveBeenCalled();
  });

  describe('with name prop', () => {
    const wrapper = shallow(<TextInput className="sample-class" name="firstName" label="First Name" type="text" />);

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
    it('should render an ErrorMessage', () => {
      const errorMessage = wrapper.find(FormGroup).find(ErrorMessage);
      expect(errorMessage.length).toBe(1);
      expect(errorMessage.prop('display')).toBe(true);
      expect(errorMessage.prop('children')).toBe('sample error');
    });
    it('should render a USWDS TextInput', () => {
      const textInput = wrapper.find(FormGroup).find(UswdsTextInput);
      expect(textInput.length).toBe(1);
      expect(textInput.prop('onChange').getMockName()).toBe('onChange');
      expect(textInput.prop('className')).toBe('sample-class');
      expect(textInput.prop('type')).toBe('text');
    });
  });

  describe('with id prop', () => {
    const wrapper = shallow(<TextInput className="sample-class" id="lastName" label="Last Name" type="text" />);

    it('should render a Label', () => {
      const label = wrapper.find(FormGroup).find(Label);
      expect(label.length).toBe(1);
      expect(label.prop('htmlFor')).toBe('lastName');
    });
    it('should render a USWDS TextInput', () => {
      const textInput = wrapper.find(FormGroup).find(UswdsTextInput);
      expect(textInput.length).toBe(1);
      expect(textInput.prop('id')).toBe('lastName');
    });
  });
  describe('with no id or name prop', () => {
    const spy = jest.spyOn(global.console, 'error');
    const wrapper = shallow(<TextInput className="sample-class" label="Some Name" type="text" />);

    it('should not render', () => {
      expect(wrapper.find(FormGroup).length).toBe(0);
    });
    it('should render console error', () => {
      expect(spy).toHaveBeenCalledWith(
        "Warning: Failed prop type: id or name required on 'TextInput'\n    in TextInput (at TextInput.test.jsx:77)",
      );
    });
  });
});

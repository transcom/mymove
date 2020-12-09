import React from 'react';
import { shallow, mount } from 'enzyme';
import { FormGroup, Label, TextInput as UswdsTextInput } from '@trussworks/react-uswds';
import { IMaskInput } from 'react-imask';

import { ErrorMessage } from '../index';

import { TextInput, TextInputMinimal } from './index';

import { TextMaskedInput } from 'components/form/fields/TextInput';

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

    it('render an ErrorMessage', () => {
      const errorMessage = wrapper.find(ErrorMessage);
      expect(errorMessage.length).toBe(1);
      expect(errorMessage.prop('display')).toBe(true);
      expect(errorMessage.prop('children')).toBe('sample error');
    });

    it('render a USWDS TextInput', () => {
      const textInput = wrapper.find(UswdsTextInput);
      expect(textInput.length).toBe(1);
      expect(textInput.prop('className')).toBe('sample-class');
      expect(textInput.prop('type')).toBe('text');
    });

    it('trigger onChange properly', () => {
      const textInput = wrapper.find(UswdsTextInput);
      expect(textInput.prop('onChange')).toBe(mockOnChange);
      textInput.simulate('change', { value: 'sample' });
      expect(mockOnChange).toHaveBeenCalledWith({ value: 'sample' });
    });
  });

  describe('with id prop', () => {
    const wrapper = shallow(<TextInputMinimal className="sample-class" id="lastName" type="text" name="lastName" />);

    it('render an ErrorMessage', () => {
      const errorMessage = wrapper.find(ErrorMessage);
      expect(errorMessage.length).toBe(1);
      expect(errorMessage.prop('display')).toBe(true);
      expect(errorMessage.prop('children')).toBe('sample error');
    });

    it('render a USWDS TextInput', () => {
      const textInput = wrapper.find(UswdsTextInput);
      expect(textInput.length).toBe(1);
      expect(textInput.prop('id')).toBe('lastName');
    });
  });

  describe('with no id or name prop', () => {
    it('render console error', () => {
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

    it('render a FormGroup', () => {
      const group = wrapper.find(FormGroup);
      expect(group.length).toBe(1);
      expect(group.prop('error')).toBe(true);
    });

    it('render a Label', () => {
      const label = wrapper.find(FormGroup).find(Label);
      expect(label.length).toBe(1);
      expect(label.prop('error')).toBe(true);
      expect(label.prop('htmlFor')).toBe('firstName');
      expect(label.prop('children')).toBe('First Name');
    });

    it('render a TextInputMinimal', () => {
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

    it('render a Label', () => {
      const label = wrapper.find(FormGroup).find(Label);
      expect(label.length).toBe(1);
      expect(label.prop('htmlFor')).toBe('lastName');
    });

    it('render a TextInputMinimal', () => {
      const textInput = wrapper.find(FormGroup).find(TextInputMinimal);
      expect(textInput.length).toBe(1);
      expect(textInput.prop('id')).toBe('lastName');
    });
  });

  describe('with no id or name prop', () => {
    it('render console error', () => {
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

describe('TextMaskedInput', () => {
  describe('with name prop', () => {
    const wrapper = shallow(
      <TextMaskedInput className="sample-class" name="firstName" label="First Name" type="text" id="firstName" />,
    );

    it('render a FormGroup', () => {
      const group = wrapper.find(FormGroup);
      expect(group.length).toBe(1);
      expect(group.prop('error')).toBe(true);
    });

    it('render a Label', () => {
      const label = wrapper.find(FormGroup).find(Label);
      expect(label.length).toBe(1);
      expect(label.prop('error')).toBe(true);
      expect(label.prop('htmlFor')).toBe('firstName');
      expect(label.prop('children')).toBe('First Name');
    });

    it('render a IMaskInput', () => {
      const textInputMinimal = wrapper.find(FormGroup).find(IMaskInput);
      expect(textInputMinimal.length).toBe(1);
      expect(textInputMinimal.prop('name')).toBe('firstName');
      expect(textInputMinimal.prop('type')).toBe('text');
      expect(textInputMinimal.prop('className')).toBe('sample-class');
    });
  });

  describe('with id prop', () => {
    const wrapper = shallow(
      <TextMaskedInput className="sample-class" id="lastName" label="Last Name" type="text" name="lastName" />,
    );

    it('render a Label', () => {
      const label = wrapper.find(FormGroup).find(Label);
      expect(label.length).toBe(1);
      expect(label.prop('htmlFor')).toBe('lastName');
    });

    it('render a IMaskInput', () => {
      const textInput = wrapper.find(FormGroup).find(IMaskInput);
      expect(textInput.length).toBe(1);
      expect(textInput.prop('id')).toBe('lastName');
    });
  });

  describe('with masking', () => {
    const wrapper = mount(
      <TextMaskedInput
        // value={"8000"}
        name="authorizedWeight"
        label="Authorized weight"
        id="authorizedWeightInput"
        mask="NUM lbs" // Nested masking imaskjs
        lazy={false} // immediate masking evaluation
        blocks={{
          // our custom masking key
          NUM: {
            mask: Number,
            thousandsSeparator: ',',
            scale: 0, // whole numbers
            signed: false, // positive numbers
          },
        }}
        value="800"
      />,
    );

    // caveat here is that the prop value will stay unmasked
    // but we can test the html element value to get the display masked value
    it('render a IMaskInput with expected unmasked prop value', () => {
      const textInput = wrapper.find(FormGroup).find(IMaskInput).getDOMNode();
      expect(textInput.value).toBe('800 lbs');
    });
  });

  describe('with no id or name prop', () => {
    it('render console error', () => {
      const spy = jest.spyOn(global.console, 'error');
      shallow(<TextMaskedInput className="sample-class" label="Some Name" type="text" />);

      expect(spy).toHaveBeenCalledWith(
        expect.stringMatching(/The prop `id` is marked as required in `TextMaskedInput`, but its value is `undefined`/),
      );
      expect(spy).toHaveBeenCalledWith(
        expect.stringMatching(
          /The prop `name` is marked as required in `TextMaskedInput`, but its value is `undefined`/,
        ),
      );
    });
  });

  afterEach(jest.resetAllMocks);
});

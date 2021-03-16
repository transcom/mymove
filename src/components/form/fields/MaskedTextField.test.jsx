import React from 'react';
import { shallow, mount } from 'enzyme';
import { FormGroup, Label } from '@trussworks/react-uswds';
import { IMaskInput } from 'react-imask';

import MaskedTextField from './MaskedTextField';

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

describe('MaskedTextField', () => {
  describe('with name prop', () => {
    const wrapper = shallow(
      <MaskedTextField className="sample-class" name="firstName" label="First Name" type="text" id="firstName" />,
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
      <MaskedTextField className="sample-class" id="lastName" label="Last Name" type="text" name="lastName" />,
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
      <MaskedTextField
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
        value="8000"
      />,
    );

    // caveat here is that the prop value will stay unmasked
    // but we can test the html element value to get the display masked value
    it('render a IMaskInput with expected unmasked prop value', () => {
      const textInput = wrapper.find(FormGroup).find(IMaskInput).getDOMNode();
      expect(textInput.value).toBe('8,000 lbs');
    });
  });

  // SKIP FOR NOW, we don't typically test component propTypes & if we are going to we should do it in a way that's more scalable
  // perhaps https://github.com/ratehub/check-prop-types
  describe.skip('with no id or name prop', () => {
    it('render console error', () => {
      const spy = jest.spyOn(global.console, 'error');
      shallow(<MaskedTextField className="sample-class" label="Some Name" type="text" />);

      expect(spy).toHaveBeenCalledWith(
        expect.stringMatching(/The prop `id` is marked as required in `MaskedTextField`, but its value is `undefined`/),
      );
      expect(spy).toHaveBeenCalledWith(
        expect.stringMatching(
          /The prop `name` is marked as required in `MaskedTextField`, but its value is `undefined`/,
        ),
      );
    });
  });

  afterEach(jest.resetAllMocks);
});

import React from 'react';
import { shallow } from 'enzyme/build';
import DayPickerInput from 'react-day-picker/DayPickerInput';

import SingleDatePicker from 'shared/JsonSchemaForm/SingleDatePicker';

describe('given a SingleDatePicker input', () => {
  describe('with empty props', () => {
    it('should render without crashing', () => {
      const wrapper = shallow(<SingleDatePicker />);
      expect(wrapper.find(DayPickerInput).length).toBe(1);
    });
  });

  describe('with filled props', () => {
    const props = {
      value: '11/11/2019',
      placeholder: 'placeholder-test',
      inputClassName: 'singledatepicker-inputclassname-test',
      format: 'DD-MMM-YY',
    };
    const wrapper = shallow(<SingleDatePicker {...props} />);

    it('should render with value', () => {
      expect(wrapper.props().value).toStrictEqual(new Date('11/11/2019'));
    });

    it('should render with format', () => {
      expect(wrapper.props().format).toBe('DD-MMM-YY');
    });

    it('should render with input prop classname', () => {
      expect(wrapper.props().inputProps.className).toBe('singledatepicker-inputclassname-test');
    });

    it('should render with placeholder', () => {
      expect(wrapper.props().placeholder).toBe('placeholder-test');
    });
  });
});

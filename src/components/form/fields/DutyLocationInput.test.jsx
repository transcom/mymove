import React from 'react';
import * as Formik from 'formik';
import { mount, shallow } from 'enzyme';
import { act } from 'react-dom/test-utils';
import AsyncSelect from 'react-select/async';
import { useField } from 'formik';

import { DutyLocationInput } from './DutyLocationInput';

import { LocationSearchBoxComponent, LocationSearchBoxContainer } from 'components/LocationSearchBox/LocationSearchBox';

const mockSetValue = jest.fn();

const useFormikContextMock = jest.spyOn(Formik, 'useFormikContext');

// Helper method

const getFieldMetaMock = () => {
  return {
    value: 'testValue',
    initialTouched: true,
    touched: false,
  };
};

// mock out formik hook as we are not testing formik
// needs to be before first describe
const metaMock = {
  touched: {},
  error: '',
  initialError: '',
  initialTouched: false,
  initialValue: '',
  value: '',
};
const fieldMock = {
  value: '',
  checked: false,
  onChange: jest.fn(),
  onBlur: jest.fn(),
  multiple: undefined,
  name: 'firstName',
};
const helperMock = {
  setValue: mockSetValue,
};

jest.mock('formik', () => ({
  ...jest.requireActual('formik'),
  useField: jest.fn(() => {
    return [
      {
        value: '',
        checked: false,
        onChange: jest.fn(),
        onBlur: jest.fn(),
        multiple: undefined,
        name: 'firstName',
      },
      {
        touched: {},
        error: '',
        initialError: '',
        initialTouched: false,
        initialValue: '',
        value: '',
      },
      {
        setValue: mockSetValue,
      },
    ];
  }),
  useFormikContext: jest.fn(() => {
    return { touched: {} };
  }),
}));

jest.mock('components/LocationSearchBox/api', () => {
  return {
    SearchDutyLocations: () =>
      new Promise((resolve) => {
        resolve([]);
      }),
    ShowAddress: () =>
      new Promise((resolve) => {
        resolve(43);
      }),
  };
});

const mockProps = { ...fieldMock, ...metaMock, ...helperMock };

beforeEach(() => {
  useFormikContextMock.mockReturnValue({
    getFieldMeta: getFieldMetaMock,
    touched: {},
  });

  useField.mockReturnValue([fieldMock, metaMock, helperMock]);
});

describe('DutyLocationInput', () => {
  describe('with all required props', () => {
    const wrapper = shallow(<DutyLocationInput {...mockProps} name="name" label="label" />);

    it('renders a Duty Location search input', () => {
      const input = wrapper.find(LocationSearchBoxContainer);
      expect(input.length).toBe(1);
    });

    it('triggers onChange properly', async () => {
      const container = wrapper.find(LocationSearchBoxContainer).dive();
      const component = container.find(LocationSearchBoxComponent).dive();
      const select = component.find(AsyncSelect);
      await select.simulate('change', { id: 1, address_id: 1 });
      expect(mockSetValue).toHaveBeenCalledWith({ address: 43, address_id: 1, id: 1 });
    });

    it('escapes regex special character input', async () => {
      metaMock.touched = false;
      const mounted = mount(<DutyLocationInput {...mockProps} name="dutyLocation" label="label" />);

      await act(async () => {
        // Only the hidden input that gets the final selected duty location has a name attribute
        mounted
          .find('input#dutyLocation-input')
          .simulate('change', { target: { id: 'dutyLocation-input', value: '-][)(*+?.\\^$|' } });
      });
      mounted.update();

      // The NoOptionsMessage component is only rendered when the 'No Options' message is displayed
      expect(mounted.exists('NoOptionsMessage')).toBe(true);
    });
  });

  afterEach(jest.resetAllMocks);
});

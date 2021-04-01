import React from 'react';
import { mount, shallow } from 'enzyme';
import { act } from 'react-dom/test-utils';
import AsyncSelect from 'react-select/async';

import { DutyStationInput } from './DutyStationInput';

import { DutyStationSearchBox } from 'scenes/ServiceMembers/DutyStationSearchBox';

const mockOnChange = jest.fn();
const mockSetValue = jest.fn();
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
      { setValue: mockSetValue },
    ],
  };
});

jest.mock('scenes/ServiceMembers/api', () => {
  return {
    ShowAddress: () => new Promise((resolve) => resolve(43)),
  };
});

describe('DutyStationInput', () => {
  describe('with all required props', () => {
    const wrapper = shallow(<DutyStationInput name="name" label="label" />);

    it('renders a Duty Station search input', () => {
      const input = wrapper.find(DutyStationSearchBox);
      expect(input.length).toBe(1);
    });

    it('triggers onChange properly', async () => {
      const input = wrapper.find(DutyStationSearchBox).dive();
      const select = input.find(AsyncSelect);
      await select.simulate('change', { id: 1, address_id: 1 });
      expect(mockSetValue).toHaveBeenCalledWith({ address: 43, address_id: 1, id: 1 });
    });

    it('escapes regex special character input', async () => {
      const mounted = mount(<DutyStationInput name="dutyStation" label="label" />);

      await act(async () => {
        // Only the hidden input that gets the final selected duty station has a name attribute
        mounted
          .find('input#dutyStation-input')
          .simulate('change', { target: { id: 'dutyStation-input', value: '-][)(*+?.\\^$|' } });
      });
      mounted.update();

      // The NoOptionsMessage component is only rendered when the 'No Options' message is displayed
      expect(mounted.exists('NoOptionsMessage')).toBe(true);
    });
  });

  afterEach(jest.resetAllMocks);
});

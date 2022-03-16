import React from 'react';
import { mount, shallow } from 'enzyme';
import { act } from 'react-dom/test-utils';
import AsyncSelect from 'react-select/async';

import { DutyLocationInput } from './DutyLocationInput';

import {
  DutyStationSearchBoxComponent,
  DutyStationSearchBoxContainer,
} from 'components/DutyStationSearchBox/DutyStationSearchBox';

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

jest.mock('components/DutyStationSearchBox/api', () => {
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

describe('DutyStationInput', () => {
  describe('with all required props', () => {
    const wrapper = shallow(<DutyLocationInput name="name" label="label" />);

    it('renders a Duty Station search input', () => {
      const input = wrapper.find(DutyStationSearchBoxContainer);
      expect(input.length).toBe(1);
    });

    it('triggers onChange properly', async () => {
      const container = wrapper.find(DutyStationSearchBoxContainer).dive();
      const component = container.find(DutyStationSearchBoxComponent).dive();
      const select = component.find(AsyncSelect);
      await select.simulate('change', { id: 1, address_id: 1 });
      expect(mockSetValue).toHaveBeenCalledWith({ address: 43, address_id: 1, id: 1 });
    });

    it('escapes regex special character input', async () => {
      const mounted = mount(<DutyLocationInput name="dutyStation" label="label" />);

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

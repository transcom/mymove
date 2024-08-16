import React from 'react';
import { mount, shallow } from 'enzyme';
import { act } from 'react-dom/test-utils';
import AsyncSelect from 'react-select/async';

import { ZipCityInput } from './ZipCityInput';

import { LocationSearchBoxComponent, LocationSearchBoxContainer } from 'components/LocationSearchBox/LocationSearchBox';

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

jest.mock('services/internalapi', () => {
  return {
    searchLocationByZipCity: () =>
      new Promise((resolve) => {
        resolve(["El Paso"]);
      }),
  };
});

const handleZipCityChange = jest.fn();

describe('ZipCityInput', () => {
  describe('with all required props', () => {
    const wrapper = shallow(
        <ZipCityInput name="zipCity"
            placeholder="Start typing a Zip Code or City..."
            label="Zip/City Lookup"
            displayAddress={false}
            handleZipCityChange={handleZipCityChange} 
        />
    );

    it('renders a Zip City search input', () => {
      const input = wrapper.find(LocationSearchBoxContainer);
      expect(input.length).toBe(1);
    });

    it('triggers onChange properly', async () => {
      const container = wrapper.find(LocationSearchBoxContainer).dive();
      const component = container.find(LocationSearchBoxComponent).dive();
      const select = component.find(AsyncSelect);
      await select.simulate('change', { city: "El Paso" });
      expect(mockSetValue).toHaveBeenCalledWith({ city: "El Paso"});
    });
  });

  afterEach(jest.resetAllMocks);
});

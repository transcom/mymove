import React from 'react';
import { shallow } from 'enzyme';
import AsyncSelect from 'react-select/async';

import { LocationInput } from './LocationInput';

import { LocationSearchBoxComponent, LocationSearchBoxContainer } from 'components/LocationSearchBox/LocationSearchBox';
import { searchLocationByZipCityState } from 'services/internalApi';

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

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  searchLocationByZipCityState: jest.fn(),
}));

const handleZipCityChange = jest.fn();

describe('ZipCityInput', () => {
  describe('with all required props', () => {
    const wrapper = shallow(
      <LocationInput
        name="zipCity"
        placeholder="Start typing a Zip Code or City..."
        label="Zip/City Lookup"
        displayAddress={false}
        handleLocationChange={handleZipCityChange}
      />,
    );

    it('renders a Zip City search input', () => {
      const input = wrapper.find(LocationSearchBoxContainer);
      expect(input.length).toBe(1);
    });

    it('triggers onChange properly', async () => {
      const cityName = 'El Paso';
      searchLocationByZipCityState.mockImplementation(() => Promise.resolve(cityName));
      const container = wrapper.find(LocationSearchBoxContainer).dive();
      const component = container.find(LocationSearchBoxComponent).dive();
      const select = component.find(AsyncSelect);
      await select.simulate('change', { city: cityName });
      expect(mockSetValue).toHaveBeenCalledWith({ city: cityName });
    });
  });

  afterEach(jest.resetAllMocks);
});

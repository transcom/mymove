import React from 'react';
import { shallow } from 'enzyme';
import { useFormikContext } from 'formik';

import RankField from './RankField';

const mockMultipleRankOptions = [
  { key: 'uuid1', value: 'PV1' },
  { key: 'uuid2', value: 'PVT2' },
];

const handleChange = jest.fn();
const mockSingleRankOption = [{ key: 'E-1', value: 'AB' }];

jest.mock('formik', () => ({
  useFormikContext: jest.fn(),
}));

describe('RankField', () => {
  beforeEach(() => {
    useFormikContext.mockReturnValue({
      setFieldValue: jest.fn(),
    });
  });

  it("should not display a dropdown if there's only one option available", () => {
    const wrapper = shallow(<RankField rankOptions={mockSingleRankOption} handleChange={handleChange} />);

    expect(wrapper.find({ 'data-testid': 'RankDropdown' }).length).toBe(0);
  });

  it('should display a dropdown if there are more than one rank option', () => {
    const wrapper = shallow(<RankField rankOptions={mockMultipleRankOptions} />);
    expect(wrapper.find({ 'data-testid': 'RankDropdown' }).length).toBe(1);
  });
});

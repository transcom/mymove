import { DimensionsField } from 'shared/JsonSchemaForm/DimensionsField';
import { shallow } from 'enzyme/build';
import React from 'react';

describe('given a dimension input', () => {
  it('should render without crashing', () => {});
  let swagger = { a: 'test', b: 2 };
  let wrapper = shallow(<DimensionsField isRequired={true} fieldName={'test'} labelText={'test'} swagger={swagger} />);
  expect(wrapper.find('.dimensions-form').length).toBe(1);
});

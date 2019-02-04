import { DimensionsField } from 'shared/JsonSchemaForm/DimensionsField';
import { shallow } from 'enzyme/build';
import React from 'react';

describe('given a dimension input', () => {
  describe('and there are no values entered', () => {
    it('should display an error', () => {});
    let swagger = { a: 'test', b: 2 };
    let wrapper = shallow(
      <DimensionsField isRequired={true} fieldName={'test'} labelText={'test'} swagger={swagger} />,
    );
    wrapper.props();
    //expect(wrapper.props().isRequired).toBe(true);
  });
});

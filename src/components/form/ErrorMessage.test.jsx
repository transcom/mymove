import React from 'react';
import { shallow } from 'enzyme';
import { ErrorMessage as UswdsErrorMessage } from '@trussworks/react-uswds';
import { ErrorMessage } from '.';

describe('ErrorMessage', () => {
  describe('with error and display true', () => {
    it('should render the USWDS ErrorMessage', () => {
      const wrapper = shallow(
        <ErrorMessage display className="sample-class">
          Error
        </ErrorMessage>,
      );

      expect(wrapper.find(UswdsErrorMessage).length).toBe(1);
      expect(wrapper.prop('className')).toBe('sample-class');
      expect(wrapper.html()).toContain('>Error<');
    });
  });
  describe('with display false', () => {
    it('should NOT render the USWDS ErrorMessage', () => {
      const wrapper = shallow(
        <ErrorMessage display={false} className="sample-class">
          Error
        </ErrorMessage>,
      );

      expect(wrapper.find(UswdsErrorMessage).length).toBe(0);
    });
  });
});

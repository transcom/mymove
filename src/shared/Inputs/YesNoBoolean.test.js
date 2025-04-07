import React from 'react';
import { render } from '@testing-library/react';

import YesNoBoolean from 'shared/Inputs/YesNoBoolean';

describe('given a single Yes No Boolean group', () => {
  const wrapper = render(<YesNoBoolean />);
  const firstBool = wrapper.getByLabelText('Yes');
  const secondBool = wrapper.getByLabelText('No');
  describe('when it loads', () => {
    it('No shoul be checked', () => {
      expect(firstBool.checked).toBeFalsy();
      expect(secondBool.checked).toBeTruthy();
    });
  });

  describe('when Yes is selected', () => {
    it('Yes should be checked', async () => {
      firstBool.click();
      expect(firstBool.checked).toBeTruthy();
    });
  });
});

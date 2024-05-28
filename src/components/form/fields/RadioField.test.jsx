import React from 'react';
import { render } from '@testing-library/react';
import { Formik } from 'formik';

import { RadioField } from './RadioField';

describe('RadioField component', () => {
  it('renders the elements that make up a field', () => {
    const { getByText, getByLabelText } = render(
      <Formik>
        <RadioField name="radioField" label="radioField" id="radioField" />
      </Formik>,
    );

    expect(getByText('radioField')).toBeInstanceOf(HTMLLabelElement);
    expect(getByLabelText('radioField')).toHaveAttribute('name', 'radioField');
    expect(getByLabelText('radioField')).toHaveAttribute('id', 'radioField');
  });

  describe('disabled', () => {
    it('disables the radio when it is disabled', () => {
      const { getByLabelText } = render(
        <Formik>
          <RadioField name="radioField" label="radioField" id="radioField" disabled />
        </Formik>,
      );

      expect(getByLabelText('radioField')).toBeDisabled();
    });
  });

  afterEach(jest.resetAllMocks);
});

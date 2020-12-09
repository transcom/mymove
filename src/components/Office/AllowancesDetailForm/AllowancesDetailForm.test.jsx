import React from 'react';
import { mount } from 'enzyme';
import { Formik } from 'formik';

import AllowancesDetailForm from './AllowancesDetailForm';

const initialValues = {
  authorizedWeight: 8000,
};

describe('AllowancesDetailForm', () => {
  const wrapper = mount(
    <Formik initalValues={initialValues}>
      <form>
        <AllowancesDetailForm />
      </form>
    </Formik>,
  );

  it('renders the Form', () => {
    expect(wrapper.find(AllowancesDetailForm).exists()).toBe(true);
  });
});

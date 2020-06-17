import React from 'react';
import { shallow } from 'enzyme';
import { Form as UswdsForm, Button } from '@trussworks/react-uswds';

import { Form } from './index';

const mockHandleReset = jest.fn();
const mockHandleSubmit = jest.fn();
// mock out formik hook as we are not testing formik
// needs to be before first describe
jest.mock('formik', () => {
  return {
    ...jest.requireActual('formik'),
    useFormikContext: () => ({
      errors: { sampleField: 'Required' },
      touched: { sampleField: true },
      handleReset: mockHandleReset,
      handleSubmit: mockHandleSubmit,
    }),
  };
});

describe('Form', () => {
  const wrapper = shallow(
    <Form className="sample-class">
      <Button type="submit">Submit</Button>
      <Button type="reset">Reset</Button>
    </Form>,
  );

  it('should render the USWDS Form', () => {
    expect(wrapper.find(UswdsForm).length).toBe(1);
  });

  it('should accept onSubmit method', () => {
    expect(wrapper.prop('onSubmit')).toBe(mockHandleSubmit);
  });

  it('should accept onReset method', () => {
    expect(wrapper.prop('onReset')).toBe(mockHandleReset);
  });

  it('should accept className', () => {
    expect(wrapper.prop('className')).toBe('sample-class');
    expect(wrapper.find(Button).length).toBe(2);
  });

  it('should call submit handler', () => {
    wrapper.simulate('submit');
    expect(mockHandleSubmit).toHaveBeenCalled();
    expect(mockHandleReset).not.toHaveBeenCalled();
  });

  it('should call reset handler', () => {
    wrapper.simulate('reset');
    expect(mockHandleSubmit).not.toHaveBeenCalled();
    expect(mockHandleReset).toHaveBeenCalled();
  });

  describe('with errorCallback', () => {
    beforeEach(() => {
      jest.spyOn(React, 'useEffect').mockImplementationOnce((f) => f());
    });

    it('passes errors to it when rendered', () => {
      const errorCallback = jest.fn();
      shallow(
        <Form errorCallback={errorCallback} className="sample-class">
          <Button type="submit">Submit</Button>
          <Button type="reset">Reset</Button>
        </Form>,
      );

      expect(errorCallback).toHaveBeenCalledWith(
        {
          sampleField: 'Required',
        },
        {
          sampleField: true,
        },
      );
    });
  });

  afterEach(jest.resetAllMocks);
});

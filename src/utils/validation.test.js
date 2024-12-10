import { scrollToViewFormikError } from 'utils/validation';

describe('scrollToViewFormikError', () => {
  let formikMock;
  let focusMock;
  beforeEach(() => {
    formikMock = { isSubmitting: true, errors: { field1: 'Error message' } };
    focusMock = jest.fn();
    document.querySelector = jest.fn().mockImplementation((selector) => {
      if (selector === '[name="field1"]' || selector === '[id="field1"]') {
        return { focus: focusMock };
      }
      return null;
    });
  });

  afterEach(() => {
    formikMock = null;
    focusMock.mockReset();
    document.querySelector.mockReset();
  });

  it('should focus on the error element when there are errors and form is submitting', () => {
    scrollToViewFormikError(formikMock);
    const errorElement = document.querySelector('[name="field1"]');
    expect(errorElement).not.toBeNull();
    expect(focusMock).toHaveBeenCalled();
  });
  it('should not focus on the error element when there are no errors', () => {
    formikMock.errors = {};
    document.querySelector = jest.fn().mockReturnValue(null);
    scrollToViewFormikError(formikMock);
    const errorElement = document.querySelector('[name="field1"]');
    expect(errorElement).toBeNull();
    expect(focusMock).not.toHaveBeenCalled();
  });
});

import React from 'react';
import * as Yup from 'yup';
import { Formik } from 'formik';
import { mount } from 'enzyme';

import { EditButton, ErrorMessage, Form, StackedTableRowForm } from './index';

describe('StackedTableRowForm', () => {
  const renderStackedTableRowForm = (submit, reset, value = 'value') => {
    return mount(
      <table className="table--stacked">
        <tbody>
          <StackedTableRowForm
            initialValues={{ fieldName: value }}
            validationSchema={Yup.object({
              ordersNumber: Yup.string().max(15, 'Must be 15 characters or less').required('Required'),
            })}
            onSubmit={submit}
            onReset={reset}
            name="fieldName"
            type="text"
            label="Field Name"
          />
        </tbody>
      </table>,
    );
  };

  describe('when NOT showing form', () => {
    it('renders a tr with correct html', () => {
      const component = renderStackedTableRowForm();
      expect(component.html()).toBe(
        '<table class="table--stacked"><tbody><tr class="stacked-table-row"><th scope="row" class="label ">Field Name</th><td><span>value</span><button type="button" class="usa-button usa-button--icon usa-button--unstyled float-right" data-testid="button"><span class="icon"><svg>edit.svg</svg></span><span>Edit</span></button></td></tr></tbody></table>',
      );
    });

    it('renders a span with nbsp when no value', () => {
      const component = renderStackedTableRowForm(null, null, null);
      expect(component.html()).toBe(
        '<table class="table--stacked"><tbody><tr class="stacked-table-row"><th scope="row" class="label ">Field Name</th><td><span>&nbsp;</span><button type="button" class="usa-button usa-button--icon usa-button--unstyled float-right" data-testid="button"><span class="icon"><svg>edit.svg</svg></span><span>Edit</span></button></td></tr></tbody></table>',
      );
    });

    it('toggles show on click of edit button', () => {
      const setShow = jest.fn();
      jest.spyOn(React, 'useState').mockReturnValueOnce([false, setShow]).mockReturnValueOnce([{}, jest.fn()]);
      const component = renderStackedTableRowForm();
      expect(component.find(EditButton).length).toBe(1);
      component.find(EditButton).simulate('click');
      expect(setShow).toHaveBeenCalledWith(true);
    });

    it('does not add error class th', () => {
      const component = renderStackedTableRowForm();
      expect(component.find('th').props().className).toBe('label ');
    });

    it('does set display on ErrorMessage', () => {
      const component = renderStackedTableRowForm();
      expect(component.find(ErrorMessage).length).toBe(1);
      expect(component.find(ErrorMessage).props().display).toBe(false);
    });
    describe('with errors', () => {
      beforeEach(() => {
        jest
          .spyOn(React, 'useState')
          .mockReturnValueOnce([false, jest.fn()])
          .mockReturnValueOnce([{ fieldName: 'Required' }, jest.fn()]);
      });

      it('adds error class th', () => {
        const component = renderStackedTableRowForm();
        expect(component.find('th').props().className).toBe('label error');
      });

      it('displays the error message and value', () => {
        const component = renderStackedTableRowForm();
        expect(component.find(ErrorMessage).length).toBe(1);
        expect(component.find(ErrorMessage).props().className).toBe('display-inline');
        expect(component.find(ErrorMessage).props().display).toBe(true);
        expect(component.find(ErrorMessage).props().children).toBe('Required');
      });
    });
  });

  describe('when showing form', () => {
    const setShow = jest.fn();
    const submit = jest.fn();
    const reset = jest.fn();

    it('renders a tr with correct html form', () => {
      jest.spyOn(React, 'useState').mockReturnValueOnce([true, setShow]).mockReturnValueOnce([{}, jest.fn()]);
      const component = renderStackedTableRowForm();
      expect(component.html()).toBe(
        '<table class="table--stacked"><tbody><tr class="stacked-table-row"><th scope="row" class="label ">Field Name</th><td><form data-testid="form" class="usa-form"><input data-testid="textInput" class="usa-input" name="fieldName" type="text" value="value"><div class="form-buttons"><button type="submit" class="usa-button" data-testid="button">Submit</button><button type="reset" class="usa-button usa-button--secondary" data-testid="button">Cancel</button></div></form></td></tr></tbody></table>',
      );
    });

    describe('event submit', () => {
      beforeEach(() => {
        jest.spyOn(React, 'useState').mockReturnValueOnce([true, setShow]).mockReturnValueOnce([{}, jest.fn()]);
      });

      it('does not show the EditButton with the form', () => {
        const component = renderStackedTableRowForm(submit, reset);
        expect(component.find(EditButton).length).toBe(0);
      });

      it('triggers submit handler', () => {
        const component = renderStackedTableRowForm(submit, reset);
        const onSubmitFunc = component.find(Formik).props().onSubmit;
        onSubmitFunc({ sample: 'submit data' });
        expect(submit).toHaveBeenCalledWith({
          sample: 'submit data',
        });
        expect(reset).not.toHaveBeenCalled();
      });

      it('resets show', () => {
        const component = renderStackedTableRowForm(submit, reset);
        const onSubmitFunc = component.find(Formik).props().onSubmit;
        onSubmitFunc({ sample: 'submit data' });
        expect(setShow).toHaveBeenCalledWith(false);
      });
    });

    describe('event reset', () => {
      beforeEach(() => {
        jest.spyOn(React, 'useState').mockReturnValueOnce([true, setShow]).mockReturnValueOnce([{}, jest.fn()]);
      });

      it('does not show the EditButton with the form', () => {
        const component = renderStackedTableRowForm(submit, reset);
        expect(component.find(EditButton).length).toBe(0);
      });

      it('triggers reset handler', () => {
        const component = renderStackedTableRowForm(submit, reset);
        const onResetFunc = component.find(Formik).props().onReset;
        onResetFunc({ sample: 'reset data' });
        expect(reset).toHaveBeenCalledWith({
          sample: 'reset data',
        });
        expect(submit).not.toHaveBeenCalled();
      });

      it('resets show', () => {
        const component = renderStackedTableRowForm(submit, reset);
        const onResetFunc = component.find(Formik).props().onReset;
        onResetFunc({ sample: 'reset data' });
        expect(setShow).toHaveBeenCalledWith(false);
      });
    });
    describe('with errors', () => {
      const setErrors = jest.fn();

      beforeEach(() => {
        jest
          .spyOn(React, 'useState')
          .mockReturnValueOnce([true, jest.fn()])
          .mockReturnValueOnce([{ fieldName: 'Required' }, setErrors]);
      });

      it('passes in an errorCallback to Form', () => {
        const component = renderStackedTableRowForm();
        expect(component.find(Form).length).toBe(1);
        component.find(Form).props().errorCallback({ fieldName: 'Required' });
        expect(setErrors).toHaveBeenCalledWith({ fieldName: 'Required' });
      });

      it('adds error class th', () => {
        const component = renderStackedTableRowForm();
        expect(component.find('th').props().className).toBe('label error');
      });
    });
  });

  afterEach(jest.clearAllMocks);
});

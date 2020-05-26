import React from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';

import { storiesOf } from '@storybook/react';
import { action } from '@storybook/addon-actions';

storiesOf('Components|Form', module).add('elements', () => (
  <div style={{ padding: '20px' }}>
    <hr />
    <h3>Form Elements</h3>
    <Formik
      initialValues={{ rejectionReason: '' }}
      validationSchema={Yup.object({
        rejectionReason: Yup.string().min(15, 'Must be 15 characters or more').required('Required'),
      })}
      onSubmit={action('Form Submit')}
      onReset={action('Form Canceled')}
    >
      <form className="usa-form">
        <label className="usa-label" htmlFor="input-type-text-example-1">
          Text input label
          <input className="usa-input" id="input-type-text-example-1" name="input-type-text-example-1" type="text" />
        </label>

        <label className="usa-label" htmlFor="input-focus">
          Text input focused
          <input className="usa-input usa-focus" id="input-focus" name="input-focus" type="text" />
        </label>

        <div className="usa-form-group usa-form-group--error">
          <label className="usa-label usa-label--error" htmlFor="input-error">
            Text input error
            <span className="usa-error-message" id="input-error-message" role="alert">
              Helpful error message
            </span>
            <input
              className="usa-input usa-input--error"
              id="input-error"
              name="input-error"
              type="text"
              aria-describedby="input-error-message"
            />
          </label>
        </div>

        <label className="usa-label" htmlFor="input-type-textarea">
          Text area label
          <textarea className="usa-textarea" id="input-type-textarea" name="input-type-textarea" />
        </label>

        <label className="usa-label" htmlFor="options">
          Dropdown label
          <select className="usa-select" name="options" id="options">
            <option value>- Select -</option>
            <option value="value1">Option A</option>
            <option value="value2">Option B</option>
            <option value="value3">Option C</option>
          </select>
        </label>

        <label className="usa-label" htmlFor="input-type-fieldset">
          Checkboxes Label
          <fieldset className="usa-fieldset" name="input-type-fieldset">
            <legend className="usa-sr-only">Historical figures 1</legend>
            <div className="usa-checkbox">
              <label className="usa-checkbox__label" htmlFor="truth">
                <input className="usa-checkbox__input" id="truth" type="checkbox" name="truth" value="truth" checked />
                Sojourner Truth
              </label>
            </div>
            <div className="usa-checkbox">
              <label className="usa-checkbox__label" htmlFor="douglass">
                <input className="usa-checkbox__input" id="douglass" type="checkbox" name="douglass" value="douglass" />
                Frederick Douglass
              </label>
            </div>
            <div className="usa-checkbox">
              <label className="usa-checkbox__label" htmlFor="carver">
                <input className="usa-checkbox__input" id="carver" type="checkbox" name="carver" disabled />
                George Washington Carver
              </label>
            </div>
          </fieldset>
        </label>

        <label className="usa-label" htmlFor="radios-fieldset">
          Radios Label
          <fieldset className="usa-fieldset" name="radios-fieldset">
            <legend className="usa-sr-only">Historical figures 2</legend>
            <div className="usa-radio">
              <label className="usa-radio__label" htmlFor="stanton">
                <input className="usa-radio__input" id="stanton" type="radio" checked name="ecs" value="stanton" />
                Elizabeth Cady Stanton
              </label>
            </div>
            <div className="usa-radio">
              <label className="usa-radio__label" htmlFor="sba">
                <input className="usa-radio__input" id="sba" type="radio" name="sba" value="anthony" />
                Susan B. Anthony
              </label>
            </div>
          </fieldset>
        </label>
      </form>
    </Formik>
  </div>
));

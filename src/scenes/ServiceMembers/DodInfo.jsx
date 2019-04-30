import { get, pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';
import { Field } from 'redux-form';
import { normalizeSSN } from 'shared/JsonSchemaForm/reduxFieldNormalizer';

import { updateServiceMember } from './ducks';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const subsetOfFields = ['affiliation', 'edipi', 'social_security_number', 'rank'];

class SSNField extends Component {
  constructor(props) {
    super(props);

    this.state = {
      focused: false,
    };

    this.localOnBlur = this.localOnBlur.bind(this);
    this.localOnFocus = this.localOnFocus.bind(this);
  }

  localOnBlur(value, something) {
    this.setState({ focused: false });
    this.props.input.onBlur(value);
  }

  localOnFocus(value, something) {
    this.setState({ focused: true });
    this.props.input.onFocus(value);
  }

  render() {
    const {
      input: { value, name },
      meta: { touched, error },
      ssnOnServer,
    } = this.props;

    let displayedValue = value;
    if (!this.state.focused && (value !== '' || ssnOnServer)) {
      displayedValue = '•••-••-••••';
    }
    const displayError = touched && error;

    // This is copied from JsonSchemaField to match the styling
    return (
      <div className={displayError ? 'usa-input-error' : 'usa-input'}>
        <label className={displayError ? 'usa-input-error-label' : 'usa-input-label'} htmlFor={name}>
          Social security number
        </label>
        {touched &&
          error && (
            <span className="usa-input-error-message" id={name + '-error'} role="alert">
              {error}
            </span>
          )}
        <input {...this.props.input} onFocus={this.localOnFocus} onBlur={this.localOnBlur} value={displayedValue} />
      </div>
    );
  }
}

const validateDodForm = (values, form) => {
  // Everything is taken care of except for SSN
  let errors = {};
  const ssn = values.social_security_number;
  const hasSSN = form.ssnOnServer;

  const validSSNPattern = RegExp('^\\d{3}-\\d{2}-\\d{4}$');
  const validSSN = validSSNPattern.test(ssn);
  const ssnPresent = ssn !== '' && ssn !== undefined;

  if (hasSSN) {
    if (ssnPresent && !validSSN) {
      errors.social_security_number = 'SSN must have 9 digits.';
    }
  } else {
    if (!ssnPresent) {
      errors.social_security_number = 'Required.';
    } else if (!validSSN) {
      errors.social_security_number = 'SSN must have 9 digits.';
    }
  }

  return errors;
};

const formName = 'service_member_dod_info';
const DodWizardForm = reduxifyWizardForm(formName, validateDodForm);

export class DodInfo extends Component {
  handleSubmit = () => {
    const pendingValues = this.props.values;
    if (pendingValues) {
      const patch = pick(pendingValues, subsetOfFields);
      return this.props.updateServiceMember(patch);
    }
  };

  render() {
    const { pages, pageKey, error, currentServiceMember, schema } = this.props;
    const initialValues = currentServiceMember ? pick(currentServiceMember, subsetOfFields) : null;

    const ssnOnServer = currentServiceMember ? currentServiceMember.has_social_security_number : false;

    return (
      <DodWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        serverError={error}
        initialValues={initialValues}
        ssnOnServer={ssnOnServer}
      >
        <h1 className="sm-heading">Create your profile</h1>
        <p>Before we can schedule your move, we need to know a little more about you.</p>
        <SwaggerField fieldName="affiliation" swagger={schema} required />
        <SwaggerField fieldName="edipi" swagger={schema} required />
        <Field name="social_security_number" component={SSNField} ssnOnServer={ssnOnServer} normalize={normalizeSSN} />
        <SwaggerField fieldName="rank" swagger={schema} required />
      </DodWizardForm>
    );
  }
}
DodInfo.propTypes = {
  schema: PropTypes.object.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ updateServiceMember }, dispatch);
}
function mapStateToProps(state) {
  const props = {
    schema: get(state, 'swaggerInternal.spec.definitions.CreateServiceMemberPayload', {}),
    values: getFormValues(formName)(state),
    ...state.serviceMember,
  };
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(DodInfo);

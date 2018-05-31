import { get, pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { updateServiceMember } from './ducks';

import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import './ServiceMembers.css';

const subsetOfFields = [
  'telephone',
  'secondary_telephone',
  'personal_email',
  'phone_is_preferred',
  'text_message_is_preferred',
  'email_is_preferred',
];

const validateContactForm = (values, form) => {
  let errors = {};

  let prefSelected = false;
  prefSelected = Boolean(
    values.phone_is_preferred ||
      values.text_message_is_preferred ||
      values.email_is_preferred,
  );
  if (!prefSelected) {
    const newError = {
      phone_is_preferred: 'Please select a preferred method of contact.',
    };
    return newError;
  }
  return errors;
};
const formName = 'service_member_contact_info';
const ContactWizardForm = reduxifyWizardForm(formName, validateContactForm);

export class ContactInfo extends Component {
  handleSubmit = () => {
    const pendingValues = this.props.formData.values;
    if (pendingValues) {
      this.props.updateServiceMember(pendingValues);
    }
  };

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      currentServiceMember,
      userEmail,
      schema,
    } = this.props;
    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = currentServiceMember
      ? pick(currentServiceMember, subsetOfFields)
      : null;
    if (initialValues && !initialValues.personal_email) {
      initialValues.personal_email = userEmail;
    }
    return (
      <ContactWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        hasSucceeded={hasSubmitSuccess}
        serverError={error}
        initialValues={initialValues}
      >
        <h1 className="sm-heading">Your Contact Info</h1>
        <SwaggerField fieldName="telephone" swagger={schema} required />
        <SwaggerField fieldName="secondary_telephone" swagger={schema} />
        <SwaggerField fieldName="personal_email" swagger={schema} required />

        <fieldset key="contact_preferences">
          <legend htmlFor="contact_preferences">
            Preferred contact method(s) during your move:
          </legend>
          <SwaggerField fieldName="phone_is_preferred" swagger={schema} />
          <SwaggerField
            fieldName="text_message_is_preferred"
            swagger={schema}
            disabled={true}
          />
          <SwaggerField fieldName="email_is_preferred" swagger={schema} />
        </fieldset>
      </ContactWizardForm>
    );
  }
}
ContactInfo.propTypes = {
  userEmail: PropTypes.string.isRequired,
  schema: PropTypes.object.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ updateServiceMember }, dispatch);
}
function mapStateToProps(state) {
  return {
    userEmail: state.user.email,
    schema: get(
      state,
      'swagger.spec.definitions.CreateServiceMemberPayload',
      {},
    ),
    formData: state.form[formName],
    ...state.serviceMember,
  };
}
export default connect(mapStateToProps, mapDispatchToProps)(ContactInfo);

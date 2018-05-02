import { pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { updateServiceMember, loadServiceMember } from './ducks';

import { reduxifyWizardForm } from 'shared/WizardPage/Form';

import './ServiceMembers.css';

const uiSchema = {
  title: 'Your Contact Info',
  order: [
    'telephone',
    'secondary_telephone',
    'personal_email',
    'contact_preferences',
  ],
  requiredFields: ['telephone', 'personal_email'],
  groups: {
    contact_preferences: {
      title: 'Preferred contact method during your move:',
      fields: [
        'phone_is_preferred',
        'text_message_is_preferred',
        'email_is_preferred',
      ],
    },
  },
};
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
  componentDidMount() {
    this.props.loadServiceMember(this.props.match.params.serviceMemberId);
  }

  handleSubmit = () => {
    const pendingValues = this.props.formData.values;
    if (pendingValues) {
      const patch = pick(pendingValues, subsetOfFields);
      this.props.updateServiceMember(patch);
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
    } = this.props;
    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = currentServiceMember
      ? pick(currentServiceMember, subsetOfFields)
      : null;
    if (initialValues && !initialValues.personal_email)
      initialValues.personal_email = userEmail;
    return (
      <ContactWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        isAsync={true}
        pageList={pages}
        pageKey={pageKey}
        hasSucceeded={hasSubmitSuccess}
        error={error}
        initialValues={initialValues}
        schema={this.props.schema}
        uiSchema={uiSchema}
      />
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
  return bindActionCreators(
    { updateServiceMember, loadServiceMember },
    dispatch,
  );
}
function mapStateToProps(state) {
  const props = {
    userEmail: state.user.email,
    schema: {},
    formData: state.form[formName],
    ...state.serviceMember,
  };
  if (state.swagger.spec) {
    props.schema = state.swagger.spec.definitions.CreateServiceMemberPayload;
  }
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(ContactInfo);

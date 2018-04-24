import { pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import {
  updateServiceMember,
  loadServiceMember,
  indexBackupContacts,
  createBackupContact,
  updateBackupContact,
} from './ducks';
import {
  renderField,
  recursivleyAnnotateRequiredFields,
} from 'shared/JsonSchemaForm';
import { reduxForm, Field, formValueSelector } from 'redux-form';
import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';

import './BackupContact.css';

const formName = 'service_member_backup_contact';
const baseForm = props => {
  const { schema, authorizeAgent } = props;
  recursivleyAnnotateRequiredFields(schema);
  const fields = schema.properties || {};

  const disableAgentPermissions = !authorizeAgent;

  return (
    <form>
      <h1 className="sm-heading">Backup Contact</h1>
      <p>
        If we can't reach you, who can we contact (such as spouse or parent)?
      </p>

      {renderField('name', fields, '')}
      {renderField('email', fields, '')}
      {renderField('telephone', fields, '')}

      <Field
        id="authorizeAgent"
        name="authorizeAgent"
        component="input"
        type="checkbox"
      />
      <label htmlFor="authorizeAgent">I authorize this person to:</label>

      <Field
        id="aaChoiceView"
        name="authorizeAgentChoice"
        component="input"
        type="radio"
        value="VIEW"
        disabled={disableAgentPermissions}
      />
      <label
        htmlFor="aaChoiceView"
        className={disableAgentPermissions ? 'disabled' : ''}
      >
        Sign for pickup or delivery in my absence, and view move details in this
        app.
      </label>

      <Field
        id="aaChoiceEdit"
        name="authorizeAgentChoice"
        component="input"
        type="radio"
        value="EDIT"
        disabled={disableAgentPermissions}
      />
      <label
        htmlFor="aaChoiceEdit"
        className={disableAgentPermissions ? 'disabled' : ''}
      >
        Represent me in all aspects of this move (this person will be invited to
        login and will be authorized with with power of attorney on your
        behalf).
      </label>
    </form>
  );
};

const validateContact = (values, form) => {
  let requiredErrors = {};
  ['name', 'email'].forEach(requiredFieldName => {
    if (
      values[requiredFieldName] === undefined ||
      values[requiredFieldName] === ''
    ) {
      requiredErrors[requiredFieldName] = 'Required.';
    }
  });
  return requiredErrors;
};

const FutureForm = reduxForm({ form: formName, validate: validateContact })(
  baseForm,
);

// in order to read the values of the from from within the form we must connect it
const selector = formValueSelector(formName);
const ConnectedFutureForm = connect(
  state => {
    const authorizeAgent = selector(state, 'authorizeAgent');
    return {
      authorizeAgent,
    };
  },
  null,
  null,
  { withRef: true },
)(FutureForm);

const addAgentChoiceInitialValues = (initialValues, permission) => {
  if (initialValues && permission) {
    if (permission == 'NONE') {
      initialValues.authorizeAgent = false;
      initialValues.authorizeAgentChoice = 'VIEW';
    } else {
      initialValues.authorizeAgent = true;
      initialValues.authorizeAgentChoice = permission;
    }
  }
};

export class BackupContact extends Component {
  componentDidMount() {
    this.props.indexBackupContacts(this.props.match.params.serviceMemberId);
  }

  handleSubmit = () => {
    const pendingValues = this.props.formData.values;
    const permission = pendingValues.authorizeAgent
      ? pendingValues.authorizeAgentChoice
      : 'NONE';
    pendingValues.permission = permission;
    console.log('SUBMITTING', pendingValues);
    if (pendingValues) {
      if (this.props.currentBackupContacts.length > 0) {
        // update existing
        const oldOne = this.props.currentBackupContacts[0];
        this.props.updateBackupContact(oldOne.id, pendingValues);
      } else {
        this.props.createBackupContact(
          this.props.match.params.serviceMemberId,
          pendingValues,
        );
      }
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
    const isValid =
      this.refs.currentForm && this.refs.currentForm.getWrappedInstance().valid;
    const isDirty =
      this.refs.currentForm && this.refs.currentForm.getWrappedInstance().dirty;
    // initialValues has to be null until there are values from the action since only the first values are taken
    var [backup1, backup2] = this.props.currentBackupContacts;
    const firstInitialValues = backup1
      ? pick(backup1, ['name', 'email', 'telephone'])
      : null;
    addAgentChoiceInitialValues(
      firstInitialValues,
      backup1 ? backup1.permission : null,
    );
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        isAsync={true}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={isValid}
        pageIsDirty={isDirty}
        hasSucceeded={hasSubmitSuccess}
        error={error}
      >
        <ConnectedFutureForm
          ref="currentForm"
          initialValues={firstInitialValues}
          handleSubmit={no_op}
          schema={this.props.schema}
        />
      </WizardPage>
    );
  }
}
BackupContact.propTypes = {
  userEmail: PropTypes.string.isRequired,
  schema: PropTypes.object.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      updateServiceMember,
      loadServiceMember,
      indexBackupContacts,
      createBackupContact,
      updateBackupContact,
    },
    dispatch,
  );
}
function mapStateToProps(state) {
  const props = {
    userEmail: state.user.email,
    currentBackupContacts: state.serviceMember.currentBackupContacts,
    loggedInUser: state.loggedInUser.loggedInUser,
    schema: {},
    formData: state.form[formName],
    ...state.serviceMember,
  };
  if (state.swagger.spec) {
    props.schema =
      state.swagger.spec.definitions.CreateServiceMemberBackupContactPayload;
  }
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(BackupContact);

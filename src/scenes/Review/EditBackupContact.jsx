import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';
import scrollToTop from 'shared/scrollToTop';

import { push } from 'react-router-redux';
import { reduxForm } from 'redux-form';

import Alert from 'shared/Alert'; // eslint-disable-line
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import { updateBackupContact } from 'scenes/ServiceMembers/ducks';
import SaveCancelButtons from './SaveCancelButtons';
import './Review.css';
import profileImage from './images/profile.png';
import { editBegin, editSuccessful, entitlementChangeBegin } from './ducks';

const editBackupContactFormName = 'edit_backup_contact';

let EditBackupContactForm = props => {
  const { schema, handleSubmit, submitting, valid } = props;
  return (
    <form onSubmit={handleSubmit}>
      <img src={profileImage} alt="" /> Backup Contact
      <hr />
      <h3 className="sm-heading">Edit Backup Contact:</h3>
      <p>Any person you assign as a backup contact must be 18 years of age or older.</p>
      <SwaggerField fieldName="name" swagger={schema} required />
      <SwaggerField fieldName="email" swagger={schema} required />
      <SwaggerField fieldName="telephone" swagger={schema} />
      <SaveCancelButtons valid={valid} submitting={submitting} />
    </form>
  );
};
EditBackupContactForm = reduxForm({
  form: editBackupContactFormName,
})(EditBackupContactForm);

class EditBackupContact extends Component {
  componentDidMount() {
    this.props.editBegin();
    this.props.entitlementChangeBegin();
    scrollToTop();
  }

  updateContact = fieldValues => {
    if (fieldValues.telephone === '') {
      fieldValues.telephone = null;
    }
    return this.props.updateBackupContact(fieldValues.id, fieldValues).then(() => {
      // This promise resolves regardless of error.
      if (!this.props.hasSubmitError) {
        this.props.editSuccessful();
        this.props.history.goBack();
      } else {
        scrollToTop();
      }
    });
  };

  render() {
    const { error, schema, backupContacts } = this.props;

    let backupContact = null;
    if (backupContacts) {
      backupContact = backupContacts[0];
    }

    return (
      <div className="usa-grid">
        {error && (
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              {error.message}
            </Alert>
          </div>
        )}
        <div className="usa-width-one-whole">
          <EditBackupContactForm initialValues={backupContact} onSubmit={this.updateContact} schema={schema} />
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    backupContacts: state.serviceMember.currentBackupContacts,
    serviceMember: state.serviceMember.currentServiceMember,
    error: get(state, 'serviceMember.error'),
    hasSubmitError: get(state, 'serviceMember.updateBackupContactError'),
    schema: get(state, 'swaggerInternal.spec.definitions.CreateServiceMemberBackupContactPayload', {}),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      push,
      updateBackupContact,
      editBegin,
      editSuccessful,
      entitlementChangeBegin,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(EditBackupContact);
